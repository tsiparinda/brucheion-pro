package repo

import (
	"brucheion/models"
	"fmt"
)

func (repo *SqlRepository) SelectUserBucketDict(userid int, urn string) (result []models.BucketDict) {

	rows, err := repo.Commands.SelectUserBucketDict.QueryContext(repo.Context, userid, urn)
	if err == nil {
		if result, err = scanDict(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBucketDict command: %v", err)
	}
	return
}

func (repo *SqlRepository) SelectUserBuckets(userid int) (values []string) {

	rows, err := repo.Commands.SelectUserBuckets.QueryContext(repo.Context, userid)
	if err == nil {
		if values, err = scanBuckets(rows); err != nil {
			repo.Logger.Panicf("scanBuckets: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBuckets command: %v", err)
	}
	return
}

func (repo *SqlRepository) SelectUserBucketKeyValue(userid int, urn string, key string) (value models.BoltJSON, err error) {

	rows, err := repo.Commands.SelectUserBucketKeyValue.QueryContext(repo.Context, userid, urn, key)
	if err == nil {
		if value.JSON, err = scanBucketKeyValue(rows); err != nil {
			err = fmt.Errorf("scanBucketKeyValue: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBucketKeyValue command: %v", err)
	}
	return
}
