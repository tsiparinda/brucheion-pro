package repo

import (
	"brucheion-pro/models"
	"brucheion/models"
	"fmt"
	"math/rand"
)

func (repo *MemoryRepo) Seed() {
	repo.catalogs = make([]models.Catalog, 3)
	var caturns, catcits, catgrps, catwrks, catvers, catexpls, onlines, languages []string
	caturns[0] = "urn:cts:sktlit:skt0001.nyaya002.J1D:"
	caturns[1] = "urn:cts:sktlit:skt0001.nyaya002.J2D:"
	caturns[2] = "urn:cts:sktlit:skt0001.nyaya002.J3D:"
	catcits[0] = "adhyāya.āhnika.sūtra"
	catcits[1] = "adhyāya.āhnika.sūtra"
	catcits[2] = "adhyāya.āhnika.sūtra"
	catgrps[0] = "Nyāya"
	catgrps[1] = "Nyāya"
	catgrps[2] = "Nyāya"
	catwrks[0] = "Nyāyabhāṣya"
	catwrks[1] = "Nyāyabhāṣya"
	catwrks[2] = "Nyāyavārttika"
	catvers[0] = "J1D"
	catvers[1] = "J2D"
	catvers[2] = "J3D"
	catexpls[0] = ""
	catexpls[1] = ""
	catexpls[2] = ""
	onlines[0] = "true"
	onlines[1] = "true"
	onlines[2] = "true"
	languages[0] = "san"
	languages[1] = "san"
	languages[2] = "san"

	for i := 0; i < 3; i++ {
		repo.catalogs[i] = models.Catalog{URN: caturns[i], Citation: catcits[i], GroupName: catgrps[i], WorkTitle: catwrks[i], VersionLabel: catvers[i], ExemplarLabel: catexpls[j], Online: onlines[j], Language: languages[i]})
	}

	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("Product_%v", i+1)
		price := rand.Float64() * float64(rand.Intn(500))
		cat := &repo.categories[rand.Intn(len(repo.categories))]
		repo.transcriptions = append(repo.transcriptions, models.Transcription{
			ID:   i + 1,
			Name: name, Price: price,
			Description: fmt.Sprintf("%v (%v)", name, cat.CategoryName),
			Category:    cat,
		})
	}
}
