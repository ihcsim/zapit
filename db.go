package urlscanner

import "io"

// Database provides a set of functionality to retrieve URL information from the database.
type Database interface {
	Exist(url string) (bool, error)
	Close() error
	Load(r io.Reader) error
}
