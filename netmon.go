package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
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

func saveToCSV(url, status string, duration time.Duration) error {
	file, err := os.OpenFile("metrics.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		time.Now().Format(time.RFC3339),
		url,
		fmt.Sprintf("%v", duration),
		status,
	}
	return writer.Write(record)
}

func monitorWithStorage(url string) {
	// Create CSV file with headers if it doesn't exist
	if _, err := os.Stat("metrics.csv"); os.IsNotExist(err) {
		file, _ := os.Create("metrics.csv")
		writer := csv.NewWriter(file)
		writer.Write([]string{"timestamp", "host", "latency", "status"})
		writer.Flush()
		file.Close()
	}

	for {
		duration, err := checkHTTP(url)
		status := "OK"
		if err != nil {
			status = "DOWN"
		}
		fmt.Printf("HTTP check %s: %v, %s\n", url, duration, status)
		if err := saveToCSV(url, status, duration); err != nil {
			fmt.Printf("Error saving to CSV: %v\n", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	go monitorWithStorage("https://www.google.com")
	select {}
}
