package db

import (
	"bufio"
	"io"
	"strings"
)

// dataSource is an in-memory map used for testing purposes.
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
	arg := strings.TrimSuffix(url, "\n")
	_, exist := dataSource[arg]
	return exist, nil
}

// Close is a no-op for an in-memory databse.
func (d *InMemoryDB) Close() error {
	return nil
}

// Load inserts data from the reader into the in-memory database.
func (d *InMemoryDB) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		dataSource[scanner.Text()] = struct{}{}
	}

	return scanner.Err()
}
