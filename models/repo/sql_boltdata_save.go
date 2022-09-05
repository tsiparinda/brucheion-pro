package repo

import (
	"brucheion/models"
	"encoding/json"
)

func (repo *SqlRepository) CreateBucketIfNotExists(bucket string, userid int) {
	result, err := repo.Commands.CreateBucketIfNotExists.ExecContext(repo.Context, userid, bucket)
	if err == nil {
		// we can use this result's methods:
		// LastInsertId() (int64, error)
		// RowsAffected() (int64, error)
		rows, _ := result.RowsAffected()
		repo.Logger.Debugf("CreateBucketIfNotExists.Count rendering rows:", rows)
		return
	} else {
		repo.Logger.Debugf("Cannot get inserted data", err.Error())
	}
}

func (repo *SqlRepository) SaveBoltData(p *models.BoltData, userid int) {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Context.Value("USER_SESSION_KEY").(int)

	repo.Logger.Debugf("SaveBoltData: userID= ", userid)
	//repo.Logger.Debugf("SaveBoltData: userID= ", repo.User.GetID())

	repo.Logger.Info("SaveBoltData are starting...")
	for i := range p.Bucket {
		newbucket := p.Bucket[i]
		catkey := p.Bucket[i]
		catvalue, _ := json.Marshal(p.Catalog[i])

		repo.CreateBucketIfNotExists(newbucket, userid)

		//put bucket data (hstore)
		result, err := repo.Commands.SaveBoltData.ExecContext(repo.Context, userid, newbucket, catkey, catvalue)
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
			result, err := repo.Commands.SaveBoltData.ExecContext(repo.Context, userid, newbucket, newkey, newvalue)
			if err == nil {
				//repo.Logger.Debugf("SaveBoltData: result of second cycle", result)
				result = result
			} else {
				repo.Logger.Debugf("SaveBoltData: Cannot get inserted data", err.Error())
			}
		}
	}
}
