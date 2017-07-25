package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	endpoint    = "/urlinfo/1"
	defaultPort = "8080"
	envHostname = "HOSTNAME"
	envPort     = "PORT"
)

func main() {
	// handle interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go catchInterrupt(quit)

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

func handleURLInfo(w http.ResponseWriter, req *http.Request) {
	log.Printf("GET %s%s", endpoint, req.URL.Path)
}

func serverURL() string {
	hostname := os.Getenv(envHostname)
	port, exist := os.LookupEnv(envPort)
	if !exist {
		port = defaultPort
	}

	return fmt.Sprintf("%s:%s", hostname, port)
}
