package repo

import (
	"brucheion/models"
	"database/sql"
	"errors"
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

func scanPassage(rows *sql.Rows) (catalog []models.Passage, err error) {
	catalog = make([]models.Passage, 0, 10)
	for rows.Next() {
		p := models.Passage{}
		err = rows.Scan(&p.Catalog, &p.FirstPassage, &p.ID, &p.ImageRefs,
			&p.LastPassage, &p.NextPassage, &p.PreviousPassage, &p.TextRefs, &p.Transcriber, &p.TranscriptionLines)
		if err == nil {
			catalog = append(catalog, p)
		} else {
			return
		}
	}
	return
}
func scanDict(rows *sql.Rows) (values []models.BucketDict, err error) {
	values = make([]models.BucketDict, 0, 10)
	for rows.Next() {
		var p models.BucketDict
		err = rows.Scan(&p.Key, &p.Value)
		if err == nil {
			values = append(values, p)
		} else {
			return
		}
	}
	return
}

func scanBuckets(rows *sql.Rows) (values []string, err error) {
	values = make([]string, 0, 10)
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			values = append(values, p)
		} else {
			return
		}
	}
	return
}

func scanBucketKeyValue(rows *sql.Rows) (value string, err error) {
	if rows == nil {
		err = errors.New("scanBucketKeyValue: this key not found")
	}
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			value = p
		} else {
			return
		}
	}
	return
}
