package models

import (
	"github.com/ThomasK81/gocite"
)

//BoltData is the container for CITE data imported from CEX files and is used in handleCEXLoad
type CiteData struct {
	Bucket  []string // workurn
	Data    []gocite.Work
	Catalog []Catalog
}

/*
from github.com/ThomasK81/gocite
// Work is a container for CTS passages that belong to the same work
type Work struct {
	WorkID      string
	Passages    []Passage
	Ordered     bool
	First, Last PassLoc
}

// Passage is the smallest CTSNode
type Passage struct {
	PassageID  string
	Range      bool
	Analysis   []Tokenisation
	Index      int
	Prev, Next PassLoc
	ImageLinks []Triple
}

// PassLoc is a container for the ID and
// the Index of a Passage for First, Last, Prev, Next
type PassLoc struct {
	Exists    bool
	PassageID string
	Index     int
}

// Tokenisation is a container for different tokenisations of the same textual information
type Tokenisation struct {
	ID            string
	Description   string
	DataStructure string
	Array         ArrayToken
}

// ArrayToken is a container for tokens represented in an array.
type ArrayToken struct {
	Type        string
	CharRepres  []string
	IntRepres   []int
	BoolRepres  []bool
	FloatRepres []float64
}

// Triple is a Simple LinkedData-Triple implementation
type Triple struct {
	Subject, Verb, Object string
}
*/
