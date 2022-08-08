package repo

import (
	"brucheion-pro/models"
	"platform/services"
)

func RegisterMemoryRepoService() {
	services.AddSingleton(func() models.Repository {
		repo := &MemoryRepo{}
		repo.Seed()
		return repo
	})
}

type MemoryRepo struct {
	transcriptions []models.Transcription
	catalogs       []models.Catalog
}

func (repo *MemoryRepo) GetTranscription(ctsurn string) (transcription models.Transcription) {
	for _, t := range repo.transcriptions {
		if t.CTSURN == ctsurn {
			transcription = t
			return
		}
	}
	return
}

func (repo *MemoryRepo) GetCatalog(urn string) (catalog models.Catalog) {
	for _, c := range repo.catalogs {
		if c.URN == urn {
			catalog = c
			return
		}
	}
	return
}
func (repo *MemoryRepo) GetTranscriptions() (results []models.Transcription) {
	return repo.transcriptions
}

func (repo *MemoryRepo) GetCatalogs() (results []models.Catalog) {
	return repo.catalogs
}

// func (repo *MemoryRepo) GetProductPage(page, pageSize int) ([]models.Product, int) {
// 	return getPage(repo.products, page, pageSize), len(repo.products)
// }

// func (repo *MemoryRepo) GetProductPageCategory(category int, page,
// 	pageSize int) (products []models.Product, totalAvailable int) {
// 	if category == 0 {
// 		return repo.GetProductPage(page, pageSize)
// 	} else {
// 		filteredProducts := make([]models.Product, 0, len(repo.products))
// 		for _, p := range repo.products {
// 			if p.Category.ID == category {
// 				filteredProducts = append(filteredProducts, p)
// 			}
// 		}
// 		return getPage(filteredProducts, page, pageSize), len(filteredProducts)
// 	}
// }

// func getPage(src []models.Product, page, pageSize int) []models.Product {
// 	start := (page - 1) * pageSize
// 	if page > 0 && len(src) > start {
// 		end := (int)(math.Min((float64)(len(src)), (float64)(start+pageSize)))
// 		return src[start:end]
// 	}
// 	return []models.Product{}
// }