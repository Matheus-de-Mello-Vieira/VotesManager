package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	url            = "http://localhost:8080"
	requestsPerSec = 1000
	duration       = 10 * time.Second
)

func main() {
	transport := &http.Transport{
		MaxIdleConnsPerHost: 1000, // Set the desired value
	}
	client := &http.Client{Transport: transport}

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Channel to control the rate of requests
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSec))
	defer ticker.Stop()

	// Timer to stop the test after the specified duration
	stopTimer := time.After(duration)

	fmt.Printf("Starting load test for %v at %d requests per second...\n", duration, requestsPerSec)

	var successCount int64
	var totalCount int64

	for {
		select {
		case <-stopTimer:
			wg.Wait()
			fmt.Println("Load test completed.")
			fmt.Printf("Sucess: %d\n", successCount)
			fmt.Printf("Total: %d\n", totalCount)
			fmt.Printf("Total: %.2f%%\n", float64(successCount)/float64(totalCount) * 100)
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				atomic.AddInt64(&totalCount, 1)
				if testAllFlow(client) {
					atomic.AddInt64(&successCount, 1)
				}
			}()
		}
	}
}
func testAllFlow(client *http.Client) bool {
	if !testHtmlPage(client, "/") {
		return false
	}

	if !testRestRequest(client, "GET", "/participants", nil) {
		return false
	}

	// if !simulateVote(client) {
	// 	return false
	// }

	if !testHtmlPage(client, "/pages/totals/rough") {
		return false
	}

	if !testRestRequest(client, "GET", "/votes/totals/rough", nil) {
		return false
	}

	return true
}

func simulateVote(client *http.Client) bool {
	body := map[string]interface{}{
		"participant_id": 1,
	}
	return testRestRequest(client, "POST", "/votes", body)
}

func testHtmlPage(client *http.Client, route string) bool {
	return testRequest(client, "GET", route, nil, "text/html")
}

func testRestRequest(client *http.Client, verb string, route string, body map[string]any) bool {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return false
	}

	return testRequest(client, verb, route, bytes.NewBuffer(jsonData), "application/json")
}

func testRequest(client *http.Client, verb string, route string, body io.Reader, contentType string) bool {
	req, err := http.NewRequest(verb, fmt.Sprint(url, route), body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return false
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
}
