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

	t.Run("200 OK", func(t *testing.T) {
		var testCases = []struct {
			path         string
			expectedBody []byte
		}{
			{path: "localhost", expectedBody: []byte(`{"URL":"localhost","IsSafe":true}`)},
			{path: "127.0.0.1", expectedBody: []byte(`{"URL":"127.0.0.1","IsSafe":true}`)},
			{path: "google.com", expectedBody: []byte(`{"URL":"google.com","IsSafe":true}`)},
			{path: "piknichok.ru", expectedBody: []byte(`{"URL":"piknichok.ru","IsSafe":false}`)},
			{path: "108.61.210.89", expectedBody: []byte(`{"URL":"108.61.210.89","IsSafe":false}`)},
			{path: "docs.google.com%3Fuser%3Drogue%26worm%3Djimbo", expectedBody: []byte(`{"URL":"docs.google.com%3Fuser%3Drogue%26worm%3Djimbo","IsSafe":false}`)},
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
				t.Errorf("Mismatch HTTP response status code. Expected %d, but got %d", http.StatusOK, actualResponse.StatusCode)
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
	})

	t.Run("400 Bad Request", func(t *testing.T) {
		path := fmt.Sprintf("%s%s", endpoint, "")
		testRequest := httptest.NewRequest("GET", path, nil)
		testResponseWriter := httptest.NewRecorder()
		handleURLInfo(testResponseWriter, testRequest)

		actualResponse := testResponseWriter.Result()
		if actualResponse.StatusCode != http.StatusBadRequest {
			t.Errorf("Mismatch HTTP response status code. Expected %d, but got %d", http.StatusBadRequest, actualResponse.StatusCode)
		}
	})
}

func TestDBHost(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		expected := fmt.Sprintf("%s:%s", defaultDBService, defaultDBPort)
		if actual := dbHost(); actual != expected {
			t.Errorf("DB host mismatch. Expected %q, but got %q", expected, actual)
		}
	})

	t.Run("From Env", func(t *testing.T) {
		if err := os.Setenv(envDBService, "my_db"); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.Unsetenv(envDBService)

		if err := os.Setenv(envDBPort, "7009"); err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		defer os.Unsetenv(envDBPort)

		expected := "my_db:7009"
		if actual := dbHost(); expected != actual {
			t.Errorf("DB host mismatch. Expected %q, but got %q", expected, actual)
		}
	})
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
