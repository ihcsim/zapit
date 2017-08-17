package zapit_test

import (
	"testing"

	"github.com/ihcsim/zapit"
	"github.com/ihcsim/zapit/internal/db"
)

func TestIsSafe(t *testing.T) {
	dbStore := &db.InMemoryDB{}
	scanner := zapit.NewScanner(dbStore)

	t.Run("Bad domains", func(t *testing.T) {
		testCases := []struct {
			url      string
			result   string
			expected bool
		}{
			{url: "google.com", result: "safe", expected: true},
			{url: "localhost:80", result: "safe", expected: true},
			{url: "localhost:8080/test", result: "safe", expected: true},
			{url: "127.0.0.1", result: "safe", expected: true},
			{url: "dont.exist.com", result: "safe", expected: true},
			{url: "linksk.us", result: "unsafe", expected: false},
			{url: "linksk.us:8080", result: "unsafe", expected: false},
			{url: "linksk.us:8080/accounts", result: "unsafe", expected: false},
			{url: "linksk.us:8080/accounts?foo=bar", result: "unsafe", expected: false},
			{url: "linksk.us:8080/accounts?foo=bar#anchor", result: "unsafe", expected: false},
			{url: "linksk.us/accounts", result: "unsafe", expected: false},
			{url: "linksk.us/accounts?foo=bar", result: "unsafe", expected: false},
			{url: "linksk.us/accounts?foo=bar#anchor", result: "unsafe", expected: false},
			{url: "linksk.us/accounts#anchor", result: "unsafe", expected: false},
			{url: "linksk.us?foo=bar", result: "unsafe", expected: false},
			{url: "linksk.us#support", result: "unsafe", expected: false},
			{url: "108.61.210.89", result: "unsafe", expected: false},
			{url: "108.61.210.89:8080", result: "unsafe", expected: false},
			{url: "108.61.210.89:8080/transactions", result: "unsafe", expected: false},
			{url: "108.61.210.89:8080/transactions?foo=bar", result: "unsafe", expected: false},
			{url: "108.61.210.89:8080/transactions?foo=bar#anchor", result: "unsafe", expected: false},
			{url: "108.61.210.89/transaction", result: "unsafe", expected: false},
			{url: "108.61.210.89/transaction?foo=bar", result: "unsafe", expected: false},
			{url: "108.61.210.89/transaction?foo=bar#anchor", result: "unsafe", expected: false},
			{url: "108.61.210.89/transaction#anchor", result: "unsafe", expected: false},
			{url: "108.61.210.89?foo=bar", result: "unsafe", expected: false},
			{url: "108.61.210.89#anchor", result: "unsafe", expected: false},
			{url: "app.linksk.us", result: "unsafe", expected: false},
			{url: "blog.linksk.us", result: "unsafe", expected: false},
			{url: "user.support.linksk.us", result: "unsafe", expected: false},
			//{url: "docs.google.com?user=rogue&worm=jimbo", result: "unsafe", expected: false},
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
}
