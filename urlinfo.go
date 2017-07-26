package urlscanner

// URLInfo contains metadata of an URL.
type URLInfo struct {
	// URL is URL in the query
	URL string

	// IsSafe is true if the URL is a malware URL.
	// Otherwise, it's false.
	IsSafe bool
}
