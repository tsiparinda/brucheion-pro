package models

//image is the container for image metadata
type image struct {
	URN      string `json:"urn"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	License  string `json:"license"`
	External bool   `json:"external"`
	Location string `json:"location"`
}
