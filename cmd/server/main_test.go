package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	urlscanner "github.com/ihcsim/url-scanner"
	"github.com/ihcsim/url-scanner/internal/db"
)

func TestHandleURLInfo(t *testing.T) {
	// set up scanner
	db := &db.InMemoryDB{}
	scanner = urlscanner.New(db)

	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	var testCases = []struct {
		path         string
		expectedBody []byte
	}{
		{path: "localhost", expectedBody: []byte(`{"url":"localhost","isSafe":true}`)},
		{path: "127.0.0.1", expectedBody: []byte(`{"url":"127.0.0.1","isSafe":true}`)},
		{path: "google.com", expectedBody: []byte(`{"url":"google.com","isSafe":true}`)},
		{path: "piknichok.ru", expectedBody: []byte(`{"url":"piknichok.ru","isSafe":false}`)},
		{path: "108.61.210.89", expectedBody: []byte(`{"url":"108.61.210.89","isSafe":false}`)},
	}

	for _, testCase := range testCases {
		path := fmt.Sprintf("%s%s", endpoint, testCase.path)
		testRequest := httptest.NewRequest("GET", path, nil)
		testResponseWriter := httptest.NewRecorder()
		handleURLInfo(testResponseWriter, testRequest)

		actualHeader := testResponseWriter.Header()
		if actual := actualHeader.Get("Content-Type"); contentType != actual {
			t.Errorf("Mismatch response content type. Expected %q, but got %q", contentType, actual)
		}

		actualResponse := testResponseWriter.Result()
		if actualResponse.StatusCode != http.StatusOK {
			t.Errorf("Mismatch HTTP response status code. Expected %q, but got %q", http.StatusOK, actualResponse.StatusCode)
		}

		actualBody := make([]byte, len(testCase.expectedBody))
		_, err := actualResponse.Body.Read(actualBody)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		if string(testCase.expectedBody) != string(actualBody) {
			t.Errorf("Mismatch respones body. Expected %s, but got %s", testCase.expectedBody, actualBody)
		}
	}
}

func TestServerURL(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		expected := ":8080"
		if actual := serverURL(); actual != expected {
			t.Errorf("Expected server to listen at %s, but got %s", expected, actual)
		}
	})

	t.Run("From Env", func(t *testing.T) {
		if err := os.Setenv(envHostname, "localhost"); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.Unsetenv(envHostname)

		if err := os.Setenv(envPort, "8088"); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.Unsetenv(envPort)

		expected := "localhost:8088"
		if actual := serverURL(); actual != expected {
			t.Errorf("Expected server to listen at %s, but got %s", expected, actual)
		}
	})
}
