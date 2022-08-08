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
