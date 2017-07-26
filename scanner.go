package urlscanner

import urlerr "github.com/ihcsim/url-scanner/internal/error"

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

// IsSafe returns an URLInfo struct containing information on the given URL.
// If the URL is a malware URL, URLInfo.IsSafe is set to true.
// Otherwise, it's set to false.
func (s *URLScanner) IsSafe(url string) (*URLInfo, error) {
	if url == "" {
		return nil, &urlerr.MalformedURLError{}
	}

	exist, err := s.db.Exist(url)
	if err != nil {
		return nil, err
	}

	return &URLInfo{URL: url, IsSafe: !exist}, nil
}
