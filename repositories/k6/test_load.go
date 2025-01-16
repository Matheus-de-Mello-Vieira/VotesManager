package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				atomic.AddInt64(&totalCount, 1)
				if simulateVote() {
					atomic.AddInt64(&successCount, 1)
				}
			}()
		}
	}
}

func simulateVote() bool {
	client := &http.Client{}

	data := map[string]interface{}{
		"participant_id": 1,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return false
	}

	req, err := http.NewRequest("POST", fmt.Sprint(url, "/votes"), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Log the response status
	fmt.Printf("Response status: %d\n", resp.StatusCode)

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
}
