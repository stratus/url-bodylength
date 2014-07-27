package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var urlsFile = flag.String("file", "", "File with URLs")

type Work struct {
	url    string
	length int
	proto  string
	err    error
}

func get(ch chan Work, u string) {
	r, err := http.Get(u)
	if err != nil {
		ch <- Work{url: u, err: err}
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ch <- Work{url: u, err: err}
		return
	}
	ch <- Work{url: u, length: len(body), proto: r.Proto}
}

func main() {
	flag.Parse()
	if *urlsFile == "" {
		log.Fatal("-file flag is required")
	}
	file, err := os.Open(*urlsFile)
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ch := make(chan Work)
	for _, url := range urls {
		fmt.Printf("Sending %s for a new worker to get.\n", url)
		go get(ch, url)
	}
	for _ = range urls {
		r := <-ch
		if r.err != nil {
			fmt.Println(r.err)
			continue
		}
		fmt.Printf("%s (%s) length %d.\n", r.url, r.proto, r.length)
	}
}
