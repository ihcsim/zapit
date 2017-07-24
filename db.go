package urlscanner

// Database provides a set of functionality to retrieve URL information from the database.
type Database interface {
	GetURLInfo(url string) (*URLInfo, error)
}
