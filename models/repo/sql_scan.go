package repo

import (
	"brucheion/models"
	"database/sql"
)

func scanBoltCatalog(rows *sql.Rows) (catalog []models.BoltCatalog, err error) {
	catalog = make([]models.BoltCatalog, 0, 10)
	for rows.Next() {
		p := models.BoltCatalog{}
		err = rows.Scan(&p.URN, &p.Citation, &p.GroupName, &p.WorkTitle,
			&p.VersionLabel, &p.ExemplarLabel, &p.Online, &p.Language)
		if err == nil {
			catalog = append(catalog, p)
		} else {
			return
		}
	}
	return
}

// func scanProduct(row *sql.Row) (p models.Product, err error) {
// 	p = models.Product{Category: &models.Category{}}
// 	err = row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category.ID,
// 		&p.Category.CategoryName)
// 	return p, err
// }
