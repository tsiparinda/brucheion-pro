package repo

import "brucheion/models"

func (repo *SqlRepository) GetBoltCatalog() (results []models.BoltCatalog) {
	rows, err := repo.Commands.GetBoltCatalog.QueryContext(repo.Context)
	if err == nil {
		if results, err = scanBoltCatalog(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec GetBoltCatalog command: %v", err)
	}
	return
}

func (repo *SqlRepository) GetPassage(userid int, urn string) (results []models.Passage) {
	rows, err := repo.Commands.GetPassage.QueryContext(repo.Context, userid, urn)
	if err == nil {
		if results, err = scanPassage(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec GetPassage command: %v", err)
	}
	return
}
