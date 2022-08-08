package models

type Repository interface {
	GetTranscription(CTSURN string) Transcription

	GetTranscriptions() []Transcription

	GetCatalog(URN string) Catalog

	GetCatalogs(URN string) Catalog

	// GetProductPage(page, pageSize int) (products []Product, totalAvailable int)

	// GetProductPageCategory(categoryId int, page, pageSize int) (products []Product,
	// 	totalAvailable int)

	// GetCategories() []Category

	Seed()
}
