package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadFromFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleFeedRequest))
	defer server.Close()

	defaultFeedURLOriginal := defaultFeedURL
	defaultFeedURL = server.URL
	defer func() {
		defaultFeedURL = defaultFeedURLOriginal
	}()

	buf := &bytes.Buffer{}
	if err := readFromFeed(buf); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	expected := "iclear.studentworkbook.pw\nlurdinha.psc.br\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("Results mismatch. Expected %q, but got %q", expected, actual)
	}
}

func TestReadFromRemoteFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleRemoteFilesRequest))
	defer server.Close()

	defaultRemoteFilesOriginal := defaultRemoteFiles
	defaultRemoteFiles = []string{server.URL}
	defer func() {
		defaultRemoteFiles = defaultRemoteFilesOriginal
	}()

	buf := &bytes.Buffer{}
	if err := readFromRemoteFiles(buf); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	expected := "afobal.cl\nalvoportas.com.br\naz-armaturen.su\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("Results mismatch. Expected %q, but got %q", expected, actual)
	}
}

func handleFeedRequest(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "testdata/rss_feed.xml")
}

func handleRemoteFilesRequest(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "testdata/static.txt")
}
