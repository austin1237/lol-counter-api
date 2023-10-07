package main

import (
	"context"
	"fmt"
	"ingest/champions"
	"ingest/dynamo"
	"ingest/source"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
)

var sourceApiUrl string
var tableName string

func init() {
	// Initialization function runs before main()
	sourceApiUrl = os.Getenv("SOURCE_API_URL")
	if sourceApiUrl == "" {
		panic("Environment variable sourceApiUrl is not set or empty.")
	}
	tableName = os.Getenv("COUNTER_TABLE_NAME")
	if tableName == "" {
		panic("Environment variable COUNTER_TABLE_NAME is not set or empty.")
	}
}

func refresh() {
	batchSize := 30
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
			go source.FetchSource(sourceApiUrl, champions.Champions[j], &wg, result)
		}

		// Wait for all goroutines in this batch to finish.
		wg.Wait()

		// Collect results from channels
		for j := i; j < end; j++ {
			data := <-result
			if data != nil {
				err := dynamo.SaveProcessedCounters(tableName, data)
				if err != nil {
					fmt.Println("failed to save", err)
				}
				fmt.Printf("Saved Champion: %s", data.Champion)
			}
		}
	}
	close(result)
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		offlineHandler()
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context) error {
	refresh()
	return nil
}

func offlineHandler() {
	refresh()
}
