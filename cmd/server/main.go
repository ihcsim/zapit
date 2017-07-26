package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"

	urlscanner "github.com/ihcsim/url-scanner"
	"github.com/ihcsim/url-scanner/internal/db"
)

const (
	endpoint    = "/urlinfo/1/"
	defaultPort = "8080"

	envHostname = "HOSTNAME"
	envPort     = "PORT"

	contentType = "application/json; charset=utf-8"
)

var scanner *urlscanner.URLScanner
var once sync.Once

func main() {
	// handle interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go catchInterrupt(quit)

	// init scanner
	db := &db.InMemoryDB{}
	initScanner(db)

	// register handler with DefaultServeMux
	http.HandleFunc(endpoint, handleURLInfo)

	// set up listener
	listenAddr := serverURL()
	log.Printf("Listening at %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Fail to start up server at %s. Cause: %s\n", listenAddr, err)
	}
}

func catchInterrupt(c <-chan os.Signal) {
	for {
		select {
		case <-c:
			log.Printf("Shutting down sever...")
			os.Exit(0)
		}
	}
}

func initScanner(db urlscanner.Database) {
	once.Do(func() {
		scanner = urlscanner.New(db)
	})
}

func handleURLInfo(w http.ResponseWriter, req *http.Request) {
	log.Printf("GET %s", req.URL.Path)

	var (
		safe bool
		err  error
	)
	url := strings.TrimPrefix(req.URL.Path, endpoint)
	if url != "" {
		safe, err = scanner.IsSafe(url)
		if err != nil {
			responseError(w, err)
			return
		}

		content := fmt.Sprintf(`{"url":"%s","isSafe":%t}`, url, safe)
		responseOK(w, []byte(content))
	}
}

func serverURL() string {
	hostname := os.Getenv(envHostname)
	port, exist := os.LookupEnv(envPort)
	if !exist {
		port = defaultPort
	}

	return fmt.Sprintf("%s:%s", hostname, port)
}

func responseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", contentType)

	content := fmt.Sprintf(`{"error": "%s"}`, err)
	w.Write([]byte(content))
}

func responseOK(w http.ResponseWriter, b []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Write(b)
}
