package zapit

import urlerr "github.com/ihcsim/zapit/internal/error"

// Scanner provides functionality to determine if the given URLs are malicious.
type Scanner struct {
	db Database
}

// NewScanner returns a new instance of a Scanner.
func NewScanner(db Database) *Scanner {
	return &Scanner{
		db: db,
	}
}

// IsSafe returns an URLInfo struct containing information on the given URL.
// If the URL is a malware URL, URLInfo.IsSafe is set to true.
// Otherwise, it's set to false.
func (s *Scanner) IsSafe(url string) (*URLInfo, error) {
	if url == "" {
		return nil, &urlerr.MalformedURLError{URL: "\"\""}
	}

	exist, err := s.db.Exist(url)
	if err != nil {
		return nil, err
	}

	return &URLInfo{URL: url, IsSafe: !exist}, nil
}
