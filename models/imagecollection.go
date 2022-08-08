package models

//imageCollection is the container for image collections along with their URN and name as strings
type imageCollection struct {
	URN        string  `json:"urn"`
	Name       string  `json:"name"`
	Collection []image `json:"images"`
}
