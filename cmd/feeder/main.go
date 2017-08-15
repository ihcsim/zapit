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

	"github.com/mmcdole/gofeed"
)

var (
	defaultFeedURL = "https://zeustracker.abuse.ch/rss.php"

	feedDataStartToken = "host="
	feedDataEndToken   = "&id="

	defaultFilesURL = []string{
		"https://zeustracker.abuse.ch/blocklist.php?download=baddomains",
		"https://zeustracker.abuse.ch/blocklist.php?download=badips",
	}
)

func main() {
	// handle errors and interrupt signal
	quit, errChan := make(chan os.Signal, 1), make(chan error)
	signal.Notify(quit, os.Interrupt)
	go handleSignal(quit, errChan)

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

}

func handleSignal(quit chan os.Signal, errChan chan error) {
	for {
		select {
		case <-quit:
			log.Println("Terminating the updates...")
			os.Exit(0)
		case err := <-errChan:
			log.Println(err)
		}
	}
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
