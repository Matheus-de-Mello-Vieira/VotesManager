package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	url            = "https://http://localhost:8080/votes"
	requestsPerSec = 1000
	duration       = 10 * time.Second // Duration of the test
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

	for {
		select {
		case <-stopTimer:
			wg.Wait()
			fmt.Println("Load test completed.")
			return
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				sendRequest()
			}()
		}
	}
}

func sendRequest() {
	client := &http.Client{}

	data := map[string]interface{}{
		"participant_id": 1,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/votes", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Log the response status
	fmt.Printf("Response status: %s\n", resp.Status)
}
