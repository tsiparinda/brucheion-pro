package models

// Old name is BoltCatalog
// Catalog contains all metadata of a CITE URN and is
//used in handleCEXLoad and page functions
type Catalog struct {
	URN           string `json:"urn"`
	Citation      string `json:"citationScheme"`
	GroupName     string `json:"groupName"`
	WorkTitle     string `json:"workTitle"`
	VersionLabel  string `json:"versionLabel"`
	ExemplarLabel string `json:"exemplarLabel"`
	Online        string `json:"online"`
	Language      string `json:"language"`
}
