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

/* from github.com/ThomasK81/gocite
// Work is a container for CTS passages that belong to the same work
type Work struct {
	WorkID      string
	Passages    []Passage
	Ordered     bool
	First, Last PassLoc
}
*/
