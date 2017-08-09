package zapit_test

import (
	"testing"

	"github.com/ihcsim/zapit"
	"github.com/ihcsim/zapit/internal/db"
	urlerr "github.com/ihcsim/zapit/internal/error"
)

func TestIsSafe(t *testing.T) {
	dbStore := &db.InMemoryDB{}
	scanner := zapit.NewScanner(dbStore)

	t.Run("well-formed URL", func(t *testing.T) {
		testCases := []struct {
			url      string
			result   string
			expected bool
		}{
			{url: "google.com", result: "safe", expected: true},
			{url: "localhost:8080", result: "safe", expected: true},
			{url: "127.0.0.1:8080", result: "safe", expected: true},
			{url: "dont.exist.com", result: "safe", expected: true},
			{url: "linksk.us", result: "unsafe", expected: false},
			{url: "piknichok.ru", result: "unsafe", expected: false},
			{url: "108.61.210.89", result: "unsafe", expected: false},
			{url: "docs.google.com?user=rogue&worm=jimbo", result: "unsafe", expected: false},
		}

		for _, testCase := range testCases {
			actual, err := scanner.IsSafe(testCase.url)
			if err != nil {
				t.Fatal("Unexpected error: ", err)
			}

			if actual.IsSafe != testCase.expected {
				t.Errorf("Expected URL %s to be %s", testCase.url, testCase.result)
			}
		}
	})

	t.Run("malformed URL", func(t *testing.T) {
		testCases := []struct {
			url string
			err error
		}{
			{url: "", err: &urlerr.MalformedURLError{}},
		}

		for _, testCase := range testCases {
			_, err := scanner.IsSafe(testCase.url)
			if !urlerr.IsMalformedURLError(err) {
				t.Fatal("Unexpected error: ", err)
			}
		}
	})
}
