package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ihcsim/zapit"
	"github.com/ihcsim/zapit/internal/db"
	"github.com/mmcdole/gofeed"
)

const (
	defaultDBService = "db"
	defaultDBPort    = "6379"

	envDBService = "DB_SERVICE"
	envDBPort    = "DB_PORT"

	dbProtocol = "tcp"
	dbTimeout  = time.Second * 2
)

var (
	defaultFeedURL = "https://zeustracker.abuse.ch/rss.php"

	feedDataStartToken = "host="
	feedDataEndToken   = "&id="

	defaultFilesURL = []string{
		"https://zeustracker.abuse.ch/blocklist.php?download=baddomains",
		"https://zeustracker.abuse.ch/blocklist.php?download=badips",
	}

	database zapit.Database
)

func main() {
	// handle errors and interrupt signal
	quit, errChan := make(chan os.Signal, 1), make(chan error)
	signal.Notify(quit, os.Interrupt)
	go handleSignal(quit, errChan)

	// set up conection to redis
	dbURL := dbHost()
	log.Printf("Connecting to database at %s", dbURL)
	if err := initDB(dbURL); err != nil {
		log.Fatal("Failed to initialize DB: ", err)
	}

	// set up buffered I/O and wait groups
	wait := sync.WaitGroup{}
	bufRSS := &bytes.Buffer{}
	bufFiles := &bytes.Buffer{}

	wait.Add(1)
	go func() {
		defer wait.Done()
		if err := readFromFeed(bufRSS); err != nil {
			errChan <- fmt.Errorf("Error encountered while reading from feed at %s: %s", defaultFeedURL, err)
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		if err := readFromFiles(bufFiles); err != nil {
			errChan <- fmt.Errorf("Error encountered while reading from files: %s", err)
		}
	}()

	// wait for all goroutine to complete, then flush the buffered writer
	wait.Wait()

	if err := database.Load(bufRSS); err != nil {
		log.Fatal("Failed to load RSS feed data into database. ", err)
	}

	if err := database.Load(bufFiles); err != nil {
		log.Fatal("Failed to load files data into database. ", err)
	}

	log.Println("Finish updating Redis database with new records")
}

func handleSignal(quit chan os.Signal, errChan chan error) {
	for {
		select {
		case <-quit:
			log.Println("Terminating feeder process...")
			if err := database.Close(); err != nil {
				log.Fatal("Failed to close database connection: ", err)
				os.Exit(1)
			}
			os.Exit(0)
		case err := <-errChan:
			log.Println(err)
		}
	}
}

func dbHost() string {
	service, exist := os.LookupEnv(envDBService)
	if !exist {
		service = defaultDBService
	}

	port, exist := os.LookupEnv(envDBPort)
	if !exist {
		port = defaultDBPort
	}

	return fmt.Sprintf("%s:%s", service, port)
}

func initDB(host string) error {
	var err error
	database, err = db.NewRedis(host, dbProtocol, dbTimeout)
	return err
}

func readFromFeed(w io.Writer) error {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(defaultFeedURL)
	if err != nil {
		return err
	}

	bw := bufio.NewWriter(w)
	defer func() {
		if err := bw.Flush(); err != nil {
			log.Fatal("Error encountered while flushing writer buffer. ", err)
		}
	}()

	for _, item := range feed.Items {
		startIndex := strings.Index(item.GUID, feedDataStartToken)
		endIndex := strings.Index(item.GUID, feedDataEndToken)
		data := item.GUID[startIndex+len(feedDataStartToken):endIndex] + "\n"

		_, err := bw.Write([]byte(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func readFromFiles(w io.Writer) error {
	var finalError error

	bw := bufio.NewWriter(w)
	defer func() {
		if err := bw.Flush(); err != nil {
			log.Fatal("Error encountered while flushing writer buffer. ", err)
		}
	}()

	for _, url := range defaultFilesURL {
		resp, err := http.Get(url)
		if err != nil {
			finalError = fmt.Errorf("%s\n%s: %s", finalError, url, err)
			continue
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			finalError = fmt.Errorf("%s\n%s: %s", finalError, url, err)
			continue
		}

		var scrubbed string
		for _, s := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(s, "#") || s == "" {
				continue
			}
			scrubbed += s + "\n"
		}

		if _, err := bw.Write([]byte(scrubbed)); err != nil {
			finalError = fmt.Errorf("%s\n%s: %s", finalError, url, err)
			continue
		}
	}

	return finalError
}

func loadIntoDB(r io.Reader) error {
	return database.Load(r)
}
