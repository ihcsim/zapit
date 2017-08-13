package db

import (
	"bytes"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

func TestRedisExist(t *testing.T) {
	mock, err := miniredis.Run()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer mock.Close()

	client, err := NewRedis(mock.Addr(), "tcp", time.Millisecond*100)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer client.Close()

	t.Run("Exist", func(t *testing.T) {
		var testCases = []struct {
			url      string
			result   string
			expected bool
		}{
			{url: "linksk.us", expected: true, result: "exist"},
			{url: "piknichok.ru", expected: true, result: "exist"},
			{url: "108.61.210.89", expected: true, result: "exist"},
			{url: "docs.google.com?user=rogue&worm=jimbo", expected: true, result: "exist"},
		}

		for _, testCase := range testCases {
			if testCase.expected {
				mock.Set(testCase.url, "")
			}
		}

		for _, testCase := range testCases {
			actual, err := client.Exist(testCase.url)
			if err != nil {
				t.Fatalf("URL %s. Unexpected error: %s", testCase.url, err)
			}

			if actual != testCase.expected {
				t.Errorf("Mismatch result. Expected %s to %s", testCase.url, testCase.result)
			}
		}
	})

	t.Run("Not Exist", func(t *testing.T) {
		var testCases = []struct {
			url      string
			result   string
			expected bool
		}{
			{url: "docs.google.com", expected: false, result: "not exist"},
			{url: "www.apple.ca", expected: false, result: "not exist"},
			{url: "localhost", expected: false, result: "not exist"},
		}

		for _, testCase := range testCases {
			exist, err := client.Exist(testCase.url)
			if err != nil {
				t.Fatalf("Unexpected error: ", err)
			}

			if exist {
				t.Errorf("Expected URL %s to not exist", testCase.url)
			}
		}
	})
}

func TestRedisLoad(t *testing.T) {
	mock, err := miniredis.Run()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer mock.Close()

	client, err := NewRedis(mock.Addr(), "tcp", time.Millisecond*100)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer client.Close()

	data := []string{"esco-gmbh.com\n", "82.221.129.19\n", "bkofchina.com\n"}
	buf := &bytes.Buffer{}
	for _, s := range data {
		if _, err := buf.Write([]byte(s)); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
	}

	if err := client.Load(buf); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	for _, s := range data {
		exist, err := client.Exist(s)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if !exist {
			t.Errorf("Expected URL %s to exist", s)
		}
	}
}
