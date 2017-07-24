package urlscanner

// dataSource is an in-memory table used for testing purposes.
// All its entries must match those in TestIsSafe()
var dataSource = map[string]*URLInfo{
	"linksk.us":     &URLInfo{url: "linksk.us"},
	"piknichok.ru":  &URLInfo{url: "piknichok.ru"},
	"108.61.210.89": &URLInfo{url: "108.61.210.89"},
}

// InMemoryDB represents an in-memory database.
type InMemoryDB struct{}

// GetURLInfo looks for the URL in the in-memory database.
// The return value of the error is always nil.
func (d *InMemoryDB) GetURLInfo(url string) (*URLInfo, error) {
	return dataSource[url], nil
}
