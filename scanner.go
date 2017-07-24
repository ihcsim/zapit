package urlscanner

var dbStore Database

// IsSafe returns true if the URL isn't a malware URL found in the server's database.
// Otherwise if the URL is in the data source, it returns false.
func IsSafe(url string) (bool, error) {
	urlInfo, err := dbStore.GetURLInfo(url)
	if err != nil {
		return false, err
	}

	if urlInfo == nil {
		return true, nil
	}
	return false, nil
}
