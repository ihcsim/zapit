package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandleURLInfo(t *testing.T) {
	t.Run("Request Path", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log.SetOutput(buf)
		flags := log.Flags()
		log.SetFlags(0)
		defer func() {
			log.SetOutput(os.Stdout)
			log.SetFlags(flags)
		}()

		path := "/localhost/account/hack"
		testRequest := httptest.NewRequest("GET", path, nil)
		handleURLInfo(nil, testRequest)

		expected := fmt.Sprintf("GET %s%s", endpoint, path)
		if actual := strings.TrimSpace(buf.String()); actual != expected {
			t.Errorf("Expected request to be %q, but got %q", expected, actual)
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
