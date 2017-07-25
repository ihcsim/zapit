package db

import urlscanner "github.com/ihcsim/url-scanner"

// dataSource is an in-memory table used for testing purposes.
// All its entries must match those in TestIsSafe()
var dataSource = map[string]*urlscanner.URLInfo{
	"linksk.us":     &urlscanner.URLInfo{URL: "linksk.us"},
	"piknichok.ru":  &urlscanner.URLInfo{URL: "piknichok.ru"},
	"108.61.210.89": &urlscanner.URLInfo{URL: "108.61.210.89"},
}

// InMemoryDB represents an in-memory database.
type InMemoryDB struct{}

// GetURLInfo looks for the URL in the in-memory database.
// The return value of the error is always nil.
func (d *InMemoryDB) GetURLInfo(url string) (*urlscanner.URLInfo, error) {
	return dataSource[url], nil
}
