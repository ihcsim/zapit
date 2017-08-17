package zapit

import (
	"net"
	"strings"
)

// domainLevel represents the depth of a domain name to scan. For example, given the URLs blog.example.com and support.eu.example.com, the scanner will check if the second-level domain name, in this case, example.com, to see if it is safe.
const domainLevel = 2

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
	exist, err := s.db.Exist(s.domain(url))
	if err != nil {
		return nil, err
	}

	return &URLInfo{URL: url, IsSafe: !exist}, nil
}

func (s *Scanner) domain(url string) string {
	indices := []int{
		strings.Index(url, ":"),
		strings.Index(url, "/"),
		strings.Index(url, "?"),
		strings.Index(url, "#"),
	}

	minIndex := len(url) - 1
	for _, index := range indices {
		i := index
		if i >= 0 && i < minIndex {
			minIndex = i
		}
	}

	// get the hostname by removing any port number, query strings and anchors
	hostname := url
	if minIndex < len(url)-1 {
		hostname = url[:minIndex]
	}
	domain := hostname

	// if this isn't a IP address, remove any subdomains
	if ip := net.ParseIP(hostname); ip == nil {
		start, count := 0, 0
		for i := len(hostname) - 1; i > 0; i-- {
			if hostname[i] == '.' {
				count += 1
			}

			if count == domainLevel {
				start = i + 1
				break
			}
		}
		domain = hostname[start:]
	}

	return domain
}
