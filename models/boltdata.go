package models

import (
	"brucheion/gocite"
)

// *** CITE Data Containers ***
//These are probably going to be retired altogether.

//BoltData is the container for CITE data imported from CEX files and is used in handleCEXLoad
type BoltData struct {
	Bucket  []string // workurn
	Data    []gocite.Work
	Catalog []BoltCatalog
}

//BoltWork is the container for BultURNs and their associated keys and is used in handleCEXLoad
type BoltWork struct {
	Key  []string // cts-node urn
	Data []BoltURN
}

//BoltURN is the container for a textpassage along with its URN, its image reference,
//and some information on preceding and anteceding works.
//Used for loading and saving CEX files, for pages, and for nodes
type BoltURN struct {
	URN      string   `json:"urn"`
	Text     string   `json:"text"`
	LineText []string `json:"linetext"`
	Previous string   `json:"previous"`
	Next     string   `json:"next"`
	First    string   `json:"first"`
	Last     string   `json:"last"`
	Index    int      `json:"sequence"`
	ImageRef []string `json:"imageref"`
}

//BoltJSON is a string representation of a JSON used in BoltRetrieve
type BoltJSON struct {
	JSON string
}

// BucketDict is the container for key-value pair from one of buckets
type BucketDict struct {
	Key   string
	Value string
}
