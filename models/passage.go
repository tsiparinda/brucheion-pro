package models

//passage is a the container for passage metadata

type Passage struct {
	ID                 string      `json:"id"`
	Transcriber        string      `json:"transcriber"`
	TranscriptionLines []string    `json:"transcriptionLines"`
	PreviousPassage    string      `json:"previousPassage"`
	NextPassage        string      `json:"nextPassage"`
	FirstPassage       string      `json:"firstPassage"`
	LastPassage        string      `json:"lastPassage"`
	ImageRefs          []string    `json:"imageRefs"`
	TextRefs           []string    `json:"textRefs"`
	Catalog            BoltCatalog `json:"catalog"`
}