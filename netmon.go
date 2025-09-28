package main

import (
	"fmt"
	"net/http"
	"time"
)

func checkHTTP(url string) (time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return time.Since(start), nil
}

func monitorHTTP(url string) {
	for {
		duration, err := checkHTTP(url)
		if err != nil {
			fmt.Printf("HTTP check %s failed: %v\n", url, err)
		} else {
			fmt.Printf("HTTP check %s: %v\n", url, duration)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	go monitorHTTP("https://www.google.com")
	select {}
}
