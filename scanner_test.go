package urlscanner_test

import (
	"testing"

	urlscanner "github.com/ihcsim/url-scanner"
	"github.com/ihcsim/url-scanner/internal/db"
)

func TestIsSafe(t *testing.T) {
	dbStore := &db.InMemoryDB{}
	urlScanner := urlscanner.New(dbStore)

	testCases := []struct {
		url      string
		result   string
		expected bool
		err      error
	}{
		{url: "google.com", result: "safe", expected: true},
		{url: "localhost:8080", result: "safe", expected: true},
		{url: "127.0.0.1:8080", result: "safe", expected: true},
		{url: "", result: "safe", expected: true},
		{url: "dont.exist.com", result: "safe", expected: true},
		{url: "linksk.us", result: "unsafe", expected: false},
		{url: "piknichok.ru", result: "unsafe", expected: false},
		{url: "108.61.210.89", result: "unsafe", expected: false},
	}

	for _, testCase := range testCases {
		actual, err := urlScanner.IsSafe(testCase.url)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if actual.IsSafe != testCase.expected {
			t.Errorf("Expected URL %s to be %s", testCase.url, testCase.result)
		}
	}
}
