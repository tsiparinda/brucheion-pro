package repo

import (
	"brucheion/models"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

func (repo *SqlRepository) CreateBucketIfNotExists(bucket string) {
	result, err := repo.Commands.CreateBucketIfNotExists.ExecContext(repo.Context, 1, bucket)
	if err == nil {
		repo.Logger.Debugf("return from saveboltdata", result)
		return
	} else {
		repo.Logger.Debugf("Cannot get inserted data", err.Error())
	}
}

func (repo *SqlRepository) SaveBoltData(p *models.BoltData) {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Context.Value("USER_SESSION_KEY").(int)
	repo.Logger.Debugf("SaveBoltData: userID= ", repo.Context)
	//services.Call(func(logger logging.Logger) { logger.Info("hello from saveboltdata") })
	repo.Logger.Debug("SaveBoltData are starting...")
	for i := range p.Bucket {
		newbucket := p.Bucket[i]
		catkey := p.Bucket[i]
		catvalue, _ := json.Marshal(p.Catalog[i])

		repo.CreateBucketIfNotExists(newbucket)

		//put bucket data (hstore)
		result, err := repo.Commands.SaveBoltData.ExecContext(repo.Context, 1, newbucket, catkey, catvalue)
		if err == nil {
			//repo.Logger.Debugf("SaveBoltData: result of first cycle", result)
			result = result
		} else {
			repo.Logger.Debugf("SaveBoltData: Cannot get inserted data", err.Error())
		}

		for j := range p.Data[i].Passages {
			newkey := p.Data[i].Passages[j].PassageID
			newvalue, _ := json.Marshal(p.Data[i].Passages[j])
			//repo.Logger.Debugf("SaveBoltData: params of second cycle", newbucket, newkey, fmt.Sprintf("%s", newvalue))
			result, err := repo.Commands.SaveBoltData.ExecContext(repo.Context, 1, newbucket, newkey, newvalue)
			if err == nil {
				//repo.Logger.Debugf("SaveBoltData: result of second cycle", result)
				result = result
			} else {
				repo.Logger.Debugf("SaveBoltData: Cannot get inserted data", err.Error())
			}
		}
	}
	user := "albatros"
	pwd, _ := os.Getwd()
	dbname := pwd + "/" + user + ".db"
	db, err := OpenBoltDB(dbname) //open bolt DB using helper function
	if err != nil {
		log.Println(fmt.Printf("handleCEXLoad: error opening userDB for writing: %s", err))
		//return err
	}
	defer db.Close()

	for i := range p.Bucket {
		newbucket := p.Bucket[i]
		/// new stuff
		//Saving the CTS Catalog data
		catkey := p.Bucket[i]
		catvalue, _ := json.Marshal(p.Catalog[i])
		catkeyb := []byte(catkey)
		catvalueb := []byte(catvalue)
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(newbucket))
			if err != nil {
				log.Fatal(err)
				//return err
			}
			return bucket.Put(catkeyb, catvalueb)
		})

		if err != nil {
			log.Fatal(err)
		}
		/// end stuff

		//saving the individual passages
		// for j := range boltdata.Data[i].Passages {
		// 	newkey := boltdata.Data[i].Passages[j].PassageID
		// 	newnode, _ := json.Marshal(boltdata.Data[i].Passages[j])
		// 	key := []byte(newkey)
		// 	value := []byte(newnode)
		// 	// store some data
		// 	err = db.Update(func(tx *bolt.Tx) error {
		// 		bucket, err := tx.CreateBucketIfNotExists([]byte(newbucket))
		// 		if err != nil {
		// 			return err
		// 		}
		// 		return bucket.Put(key, value)
		// 	})

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }
	}

	log.Println("CEX data loaded into Brucheion successfully!!!.")
	//return nil

	// if p.ID == 0 {
	// 	result, err := repo.Commands.SaveBoltCatalog.ExecContext(repo.Context, p.Name,
	// 		p.Description, p.Category.ID, p.Price)
	// 	if err == nil {
	// 		id, err := result.LastInsertId()
	// 		if err == nil {
	// 			p.ID = int(id)
	// 			return
	// 		} else {
	// 			repo.Logger.Panicf("Cannot get inserted ID: %v", err.Error())
	// 		}
	// 	} else {
	// 		repo.Logger.Panicf("Cannot exec SaveProduct command: %v", err.Error())
	// 	}
	// } else {
	// 	result, err := repo.Commands.UpdateProduct.ExecContext(repo.Context, p.Name,
	// 		p.Description, p.Category.ID, p.Price, p.ID)
	// 	if err == nil {
	// 		affected, err := result.RowsAffected()
	// 		if err == nil && affected != 1 {
	// 			repo.Logger.Panicf("Got unexpected rows affected: %v", affected)
	// 		} else if err != nil {
	// 			repo.Logger.Panicf("Cannot get rows affected: %v", err)
	// 		}
	// 	} else {
	// 		repo.Logger.Panicf("Cannot exec Update command: %v", err.Error())
	// 	}
	// }

}
