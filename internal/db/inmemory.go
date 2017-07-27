package db

// dataSource is an in-memory table used for testing purposes.
var dataSource = map[string]struct{}{
	"linksk.us":                             struct{}{},
	"piknichok.ru":                          struct{}{},
	"108.61.210.89":                         struct{}{},
	"docs.google.com?user=rogue&worm=jimbo": struct{}{},
}

// InMemoryDB represents an in-memory database.
type InMemoryDB struct{}

// Exist looks for the URL in the in-memory database.
// The return value of the error is always nil.
func (d *InMemoryDB) Exist(url string) (bool, error) {
	_, exist := dataSource[url]
	return exist, nil
}
