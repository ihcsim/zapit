package urlscanner

// URLScanner provides functionality to determine if the given URLs are malicious.
type URLScanner struct {
	db Database
}

// New returns a new instance of a URLScanner.
func New(db Database) *URLScanner {
	return &URLScanner{
		db: db,
	}
}

// IsSafe returns true if the URL isn't a malware URL found in the server's database.
// Otherwise if the URL is in the data source, it returns false.
func (s *URLScanner) IsSafe(url string) (bool, error) {
	urlInfo, err := s.db.GetURLInfo(url)
	if err != nil {
		return false, err
	}

	if urlInfo == nil {
		return true, nil
	}
	return false, nil
}
