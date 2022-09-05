package api

import (
	"brucheion/gocite"
	"brucheion/models"
	"brucheion/models/repo"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

type CEXuploadHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

func (h CEXuploadHandler) PostCEXupload(data string) actionresults.ActionResult {

	h.Logger.Debugf("api.ImportCEX: len(data)= ", len(data))

	err := h.importCEX(data)

	if err != nil {
		return &StatusCodeResult{http.StatusBadRequest}
	}

	return &StatusCodeResult{http.StatusOK}
}

func (h CEXuploadHandler) importCEX(data string) error {
	userid := h.User.GetID()
	h.Logger.Debugf("api.importCEX: userID= ", userid)
	var urns, areas []string
	var catalog []models.BoltCatalog

	//read in the relations of the CEX file cutting away all unnecessary signs
	if strings.Contains(data, "#!relations") {
		relations := strings.Split(data, "#!relations")[1]
		relations = strings.Split(relations, "#!")[0]
		re := regexp.MustCompile("(?m)[\r\n]*^//.*$")
		relations = re.ReplaceAllString(relations, "")

		reader := csv.NewReader(strings.NewReader(relations))
		reader.Comma = '#'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = 3

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			if strings.Contains(line[1], "appearsOn") {
				urns = append(urns, line[0])
				areas = append(areas, line[2])
			}
		}
	}
	h.Logger.Debug("ImportCEX.1")
	//read in the ctscatalog (if exists)
	if strings.Contains(data, "#!ctscatalog") {
		ctsCatalog := strings.Split(data, "#!ctscatalog")[1]
		ctsCatalog = strings.Split(ctsCatalog, "#!")[0]
		re := regexp.MustCompile("(?m)[\r\n]*^//.*$")
		ctsCatalog = re.ReplaceAllString(ctsCatalog, "")

		var caturns, catcits, catgrps, catwrks, catvers, catexpls, onlines, languages []string

		reader := csv.NewReader(strings.NewReader(ctsCatalog))
		reader.Comma = '#'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		reader.TrimLeadingSpace = true

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				h.Logger.Panicf("api.importCEX: Error reading file", error)
			}

			switch {
			case len(line) == 8:
				if line[0] != "urn" {
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
				h.Logger.Debug("api.importCEX: Catalogue Data not well formatted")
			}
		}
		for j := range caturns {
			catalog = append(catalog, models.BoltCatalog{URN: caturns[j], Citation: catcits[j], GroupName: catgrps[j], WorkTitle: catwrks[j], VersionLabel: catvers[j], ExemplarLabel: catexpls[j], Online: onlines[j], Language: languages[j]})
		}
	}
	h.Logger.Debug("api.importCEX:2")
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
			//fmt.Println(line)
			h.Logger.Panicf("api.importCEX: Error reading file2", error)
		}
		switch {
		case len(line) == 2:
			texturns = append(texturns, line[0])
			text = append(text, line[1])
		case len(line) > 2:
			texturns = append(texturns, line[0])
			var textstring string
			for j := 1; j < len(line); j++ {
				textstring = textstring + line[j]
			}
			text = append(text, textstring)
		case len(line) < 2:
			h.Logger.Panicf("api.importCEX: Wrong line", line)
		}
	}

	works := append([]string(nil), texturns...)
	for i := range texturns {
		works[i] = strings.Join(strings.Split(texturns[i], ":")[0:4], ":") + ":"
	}
	works = repo.RemoveDuplicatesUnordered(works)
	var boltworks []gocite.Work
	var sortedcatalog []models.BoltCatalog
	for i := range works {
		work := works[i]
		testexist := false
		for j := range catalog {
			if catalog[j].URN == work {
				sortedcatalog = append(sortedcatalog, catalog[j])
				testexist = true
			}
		}
		if testexist == false {
			h.Logger.Debugf("api.importCEX: urn has no catalog entry", works[i])
			sortedcatalog = append(sortedcatalog, models.BoltCatalog{})
		}

		var passages []gocite.Passage
		for j := range texturns {
			if strings.Contains(texturns[j], work) {
				var textareas []gocite.Triple
				if repo.Contains(urns, texturns[j]) {
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
		boltworks = append(boltworks, gocite.Work{WorkID: work, Passages: passages, Ordered: true, First: gocite.PassLoc{Exists: true, PassageID: passages[0].PassageID, Index: 0}, Last: gocite.PassLoc{Exists: true, PassageID: passages[len(passages)-1].PassageID, Index: len(passages) - 1}})
	}
	boltdata := models.BoltData{Bucket: works, Data: boltworks, Catalog: sortedcatalog}

	h.Logger.Debug("api.importCEX:3")
	h.Repository.SaveBoltData(&boltdata, userid)

	h.Logger.Info("api.importCEX: CEX data loaded into Brucheion successfully.!")
	return nil
}
