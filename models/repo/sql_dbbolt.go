package repo

import (
	"time"

	"github.com/boltdb/bolt"
)

//openBoltDB returns an opened Bolt Database for given dbName.
func OpenBoltDB(dbName string) (*bolt.DB, error) {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 30 * time.Second}) //open DB with - wr- --- ---
	if err != nil {
		return nil, err
	}
	return db, nil
}
