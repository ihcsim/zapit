package db

import (
	"bytes"
	"testing"
)

func TestInMemoryDBExist(t *testing.T) {
	var testCases = []struct {
		url      string
		result   string
		expected bool
	}{
		{url: "linksk.us", expected: true, result: "exist"},
		{url: "piknichok.ru", expected: true, result: "exist"},
		{url: "108.61.210.89", expected: true, result: "exist"},
		{url: "docs.google.com?user=rogue&worm=jimbo", expected: true, result: "exist"},
		{url: "docs.google.com", expected: false, result: "not exist"},
		{url: "www.apple.ca", expected: false, result: "not exist"},
		{url: "localhost", expected: false, result: "not exist"},
	}

	db := &InMemoryDB{}
	for _, testCase := range testCases {
		actual, err := db.Exist(testCase.url)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if actual != testCase.expected {
			t.Errorf("Mismatch result. Expected %s to %s", testCase.url, testCase.result)
		}
	}
}

func TestInMemoryDBLoad(t *testing.T) {
	data := []string{"esco-gmbh.com\n", "82.221.129.19\n", "bkofchina.com\n"}
	buf := &bytes.Buffer{}
	for _, s := range data {
		if _, err := buf.Write([]byte(s)); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	}

	db := &InMemoryDB{}
	if err := db.Load(buf); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	for _, s := range data {
		exist, err := db.Exist(s)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if !exist {
			t.Errorf("Expected URL %s to exist in the database", s)
		}
	}
}
