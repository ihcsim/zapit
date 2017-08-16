package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
)

func TestUpgradeInterval(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		expected := defaultDBUpdateInterval
		actual, err := updateInterval()
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		if expected != actual {
			t.Error("Upgrade interval mismatch. Expected %s, but got %s", expected, actual)
		}
	})

	t.Run("Env vars", func(t *testing.T) {
		if err := os.Setenv(envDBUpdateInterval, "20m"); err != nil {
			t.Fatal("Unexpected error", err)
		}
		defer os.Unsetenv(envDBService)

		expected, err := time.ParseDuration("20m")
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		actual, err := updateInterval()
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		if actual != expected {
			t.Errorf("Upgrade interval mismatch. Expected %s but got %s", expected, actual)
		}
	})
}

func TestDBHost(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		expected := "db:6379"
		if actual := dbHost(); expected != actual {
			t.Errorf("DB host mismatch. Expected %q, but got %q", expected, actual)
		}
	})

	t.Run("Env vars", func(t *testing.T) {
		if err := os.Setenv(envDBService, "my-redis"); err != nil {
			t.Fatal("Unexpected error", err)
		}
		defer os.Unsetenv(envDBService)

		if err := os.Setenv(envDBPort, "7000"); err != nil {
			t.Fatal("Unexpected error", err)
		}
		defer os.Unsetenv(envDBPort)

		expected := "my-redis:7000"
		if actual := dbHost(); expected != actual {
			t.Errorf("DB host mismatch. Expected %q, but got %q", expected, actual)
		}
	})
}

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

	defaultRemoteFilesOriginal := defaultFilesURL
	defaultFilesURL = []string{server.URL}
	defer func() {
		defaultFilesURL = defaultRemoteFilesOriginal
	}()

	buf := &bytes.Buffer{}
	if err := readFromFiles(buf); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	expected := "afobal.cl\nalvoportas.com.br\naz-armaturen.su\n"
	if actual := buf.String(); actual != expected {
		t.Errorf("Results mismatch. Expected %q, but got %q", expected, actual)
	}
}

func TestLoadIntoDB(t *testing.T) {
	mock, err := miniredis.Run()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer mock.Close()

	if err := initDB(mock.Addr()); err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer database.Close()

	testData := []string{"bright.su\n", "l3d1.pp.ru\n", "asscomminc.tk\n"}
	b := &bytes.Buffer{}
	for _, data := range testData {
		b.WriteString(data)
	}
	if err := loadIntoDB(b); err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	// make sure test data exist in database
	for _, data := range testData {
		exist, err := database.Exist(data)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}

		if !exist {
			t.Errorf("Expected URL %q to exist in the database", data)
		}
	}
}

func handleFeedRequest(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "testdata/rss_feed.xml")
}

func handleRemoteFilesRequest(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "testdata/static.txt")
}
