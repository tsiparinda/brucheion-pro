package models

type Repository interface {

	// GetProduct(id int) Product
	// GetProducts() []Product
	// SaveProduct(*Product)

	// GetProductPage(page, pageSize int) (products []Product, totalAvailable int)

	// GetProductPageCategory(categoryId int, page, pageSize int) (products []Product,
	//     totalAvailable int)

	// GetCategories() []Category
	// SaveCategory(*Category)

	// GetOrder(id int) Order
	// GetOrders() []Order
	// SaveOrder(*Order)
	// SetOrderShipped(*Order)
	//SaveCEX(*BoltData)
	//GetBoltData() []BoltData
	// Seed()
	// Init()
	SelectUserBucketDict(userid int, urn string) []BucketDict
	SelectUserBuckets(userid int) []string
	SelectUserBucketKeyValue(userid int, urn string, key string) (BoltJSON, error)
	GetPassage(userid int, urn string) []Passage
	//BoltRetrieve(userid int, urn string, key string)
	SaveBoltData(*BoltData, int)
	CreateBucketIfNotExists(bucket string, userid int)
	GetBoltCatalog() []BoltCatalog
	LoadMigrations()
}
