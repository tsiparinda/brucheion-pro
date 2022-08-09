package repo

import (
	"brucheion-pro/models"
	"brucheion-pro/moving"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ThomasK81/gocite"
	"github.com/boltdb/bolt"
)

func loadCEX(data string, user string) error {
	var urns, areas []string // txt arrays for urn and area in this file

	//read in the relations of the CEX file cutting away all unnecessary signs
	if strings.Contains(data, "#!relations") {
		relations := strings.Split(data, "#!relations")[1]
		relations = strings.Split(relations, "#!")[0]
		re := regexp.MustCompile("(?m)[\r\n]*^//.*$")
		relations = re.ReplaceAllString(relations, "")

		reader := csv.NewReader(strings.NewReader(relations))
		reader.Comma = '#'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = 3 // #!relations contains 3 fields

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				panic(error)
			}
			if strings.Contains(line[1], "appearsOn") {
				urns = append(urns, line[0])   // fill out the text array urns only!!!
				areas = append(areas, line[2]) // fill out the text array areas only!!!
			}
		}
	}

	var catalog []models.Catalog

	//read in the ctscatalog (if exists)
	if strings.Contains(data, "#!ctscatalog") { // 8 fields
		ctsCatalog := strings.Split(data, "#!ctscatalog")[1]
		ctsCatalog = strings.Split(ctsCatalog, "#!")[0]
		re := regexp.MustCompile("(?m)[\r\n]*^//.*$")
		ctsCatalog = re.ReplaceAllString(ctsCatalog, "")

		var caturns, catcits, catgrps, catwrks, catvers, catexpls, onlines, languages []string
		// var languages [][]string

		reader := csv.NewReader(strings.NewReader(ctsCatalog))
		reader.Comma = '#'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1 // variable?
		reader.TrimLeadingSpace = true

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				panic(error)
			}

			switch {
			case len(line) == 8:
				if line[0] != "urn" { // must be urn::::::
					caturns = append(caturns, line[0])
					catcits = append(catcits, line[1])
					catgrps = append(catgrps, line[2])
					catwrks = append(catwrks, line[3])
					catvers = append(catvers, line[4])
					catexpls = append(catexpls, line[5])
					onlines = append(onlines, line[6])
					languages = append(languages, line[7])
				}
			case len(line) != 8:
				log.Println("Catalogue Data not well formatted")
			}
		}
		for j := range caturns {
			catalog = append(catalog, models.Catalog{URN: caturns[j], Citation: catcits[j], GroupName: catgrps[j], WorkTitle: catwrks[j], VersionLabel: catvers[j], ExemplarLabel: catexpls[j], Online: onlines[j], Language: languages[j]})
		}
	}

	//read in the cts data
	ctsdata := strings.Split(data, "#!ctsdata")[1]
	ctsdata = strings.Split(ctsdata, "#!")[0]
	re := regexp.MustCompile("(?m)[\r\n]*^//.*$")
	ctsdata = re.ReplaceAllString(ctsdata, "")

	reader := csv.NewReader(strings.NewReader(ctsdata))
	reader.Comma = '#'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	var texturns, text []string

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println(line)
			log.Fatal(error)
		}
		switch {
		case len(line) == 2:
			texturns = append(texturns, line[0]) // urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.1
			text = append(text, line[1])
		case len(line) > 2: // join all fields after 0 to 1
			texturns = append(texturns, line[0])
			var textstring string
			for j := 1; j < len(line); j++ {
				textstring = textstring + line[j]
			}
			text = append(text, textstring)
		case len(line) < 2:
			fmt.Println("Wrong line:", line)
		}
	}

	works := append([]string(nil), texturns...) // urn:cts:sktlit:skt0001.nyaya002 - only 5 (4?) fields, unque
	for i := range texturns {
		works[i] = strings.Join(strings.Split(texturns[i], ":")[0:4], ":") + ":"
	}
	works = moving.RemoveDuplicatesUnordered(works)
	var boltworks []gocite.Work // passages for URN
	var sortedcatalog []models.Catalog
	for i := range works {
		work := works[i]
		testexist := false
		for j := range catalog {
			if catalog[j].URN == work {
				sortedcatalog = append(sortedcatalog, catalog[j]) // ordered catalog
				testexist = true                                  // and check if cts data exists
			}
		}
		if testexist == false {
			log.Println(works[i], " has no catalog entry")
			sortedcatalog = append(sortedcatalog, models.Catalog{})
		}

		var passages []gocite.Passage
		for j := range texturns {
			if strings.Contains(texturns[j], work) {
				var textareas []gocite.Triple
				if moving.Contains(urns, texturns[j]) {
					for k := range urns {
						if urns[k] == texturns[j] {
							textareas = append(textareas, gocite.Triple{Subject: texturns[j],
								Verb:   "urn:cite2:dse:verbs.v1:appears_on",
								Object: areas[k]})
						}
					}
				}
				linetext := strings.Replace(text[j], "-NEWLINE-", "\r\n", -1)
				passages = append(passages, gocite.Passage{PassageID: texturns[j],
					Range: false,
					Text: gocite.EncText{Brucheion: text[j],
						TXT: linetext},
					ImageLinks: textareas})
			}
		}
		//assign Next and Prev fields for all passages
		for j := range passages {
			passages[j].Index = j
			switch {
			case j+1 == len(passages):
				passages[j].Next = gocite.PassLoc{Exists: false}
			default:
				passages[j].Next = gocite.PassLoc{Exists: true, PassageID: passages[j+1].PassageID, Index: j + 1}
			}
			switch {
			case j == 0:
				passages[j].Prev = gocite.PassLoc{Exists: false}
			default:
				passages[j].Prev = gocite.PassLoc{Exists: true, PassageID: passages[j-1].PassageID, Index: j - 1}
			}
		}
		/*workToBeSaved, _ := gocite.SortPassages(gocite.Work{WorkID: work, Passages: passages, Ordered: true, First: gocite.PassLoc{Exists: true, PassageID: passages[0].PassageID, Index: 0}, Last: gocite.PassLoc{Exists: true, PassageID: passages[len(passages)-1].PassageID, Index: len(passages) - 1}})
		boltworks = append(boltworks, workToBeSaved)*/
		boltworks = append(boltworks, gocite.Work{WorkID: work, Passages: passages, Ordered: true, First: gocite.PassLoc{Exists: true, PassageID: passages[0].PassageID, Index: 0}, Last: gocite.PassLoc{Exists: true, PassageID: passages[len(passages)-1].PassageID, Index: len(passages) - 1}})
	}
	citedata := models.CiteData{Bucket: works, Data: boltworks, Catalog: sortedcatalog}
	// urn:cts:sktlit:skt0001.nyaya002 = BUCKET
	// urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.1 =
	// ???

	// write to database
	pwd, _ := os.Getwd()
	dbname := pwd + "/" + user + ".db"
	db, err := moving.OpenBoltDB(dbname) //open bolt DB using helper function
	if err != nil {
		log.Println(fmt.Printf("handleCEXLoad: error opening userDB for writing: %s", err))
		return err
	}
	defer db.Close()
	for i := range citedata.Bucket {
		newbucket := citedata.Bucket[i]
		newcatkey := citedata.Bucket[i]
		newcatvalue, _ := json.Marshal(citedata.Catalog[i])
		catkey := []byte(newcatkey)
		catvalue := []byte(newcatvalue)
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(newbucket)) /// new stuff
			if err != nil {
				return err
			}
			return bucket.Put(catkey, catvalue) //Saving the CTS Catalog data
		})

		if err != nil {
			log.Fatal(err)
		}
		/// end stuff

		//saving the individual passages
		for j := range citedata.Data[i].Passages {
			newkey := citedata.Data[i].Passages[j].PassageID
			newnode, _ := json.Marshal(citedata.Data[i].Passages[j])
			key := []byte(newkey)
			value := []byte(newnode)
			// store some data
			err = db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte(newbucket))
				if err != nil {
					return err
				}
				return bucket.Put(key, value)
			})

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("CEX data loaded into Brucheion successfully.")
	return nil
}
