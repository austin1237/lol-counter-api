package main

import (
	"fmt"
	"ingest/champions"
	"ingest/dynamo"
	"ingest/source"
	"net/url"
	"os"
	"sync"
)

var sourceApiUrl string

func init() {
	// Initialization function runs before main()
	sourceApiUrl = os.Getenv("SOURCE_API_URL")
	if sourceApiUrl == "" {
		panic("Environment variable sourceApiUrl is not set or empty.")
	}
}

func main() {
	batchSize := 16
	totalURLs := len(champions.Champions)
	var wg sync.WaitGroup
	result := make(chan *source.ProcessedCounters, batchSize)

	for i := 0; i < totalURLs; i += batchSize {
		end := i + batchSize
		if end > totalURLs {
			end = totalURLs
		}

		// Launch goroutines to fetch sources concurrently.
		for j := i; j < end; j++ {
			wg.Add(1)
			championName := url.QueryEscape(champions.Champions[j])
			go source.FetchSource(sourceApiUrl, championName, &wg, result)
		}

		// Wait for all goroutines in this batch to finish.
		wg.Wait()

		// Collect results from channels
		for j := i; j < end; j++ {
			data := <-result
			if data != nil {
				err := dynamo.SaveProcessedCounters(data)
				if err != nil {
					fmt.Println("failed to save", err)
				}
				fmt.Printf("Champion: %s, Win Rate: %s\n", data.Champion, data.Counters)
			}
		}
	}
	close(result)
}
