package main

import (
	"context"
	"ingest/champions"
	"ingest/dynamo"
	"ingest/source"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog/log"
)

var sourceApiUrl string
var tableName string
var batchSize int
var dbInterface *dynamodb.DynamoDB

func init() {
	// Initialization function runs before main()
	sourceApiUrl = os.Getenv("COUNTER_SOURCE_API_URL")
	if sourceApiUrl == "" {
		panic("Environment variable COUNTER_SOURCE_API_URL is not set or empty.")
	}
	tableName = os.Getenv("COUNTER_TABLE_NAME")
	if tableName == "" {
		panic("Environment variable COUNTER_TABLE_NAME is not set or empty.")
	}
	// Check if BATCH_SIZE is set
	batchSizeEnv := os.Getenv("BATCH_SIZE")
	if batchSizeEnv != "" {
		// Attempt to convert the environment variable to an int
		if val, err := strconv.Atoi(batchSizeEnv); err == nil {
			// Successfully converted to an int
			batchSize = val
		} else {
			// Handle the error (e.g., log it, set a different default, etc.)
			log.Error().Err(err).Msg("Error converting BATCH_SIZE to an integer")
			batchSize = 1
		}
	} else {
		// BATCH_SIZE is not set, so use the default value of 1
		batchSize = 1
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Error().Err(err).Msg("Error setting up dynamo")
	}

	// Create a new DynamoDB client
	dbInterface = dynamodb.New(sess)
}

func refresh(champions []string) {
	totalURLs := len(champions)

	var wg sync.WaitGroup
	result := make(chan *source.ProcessedCounters, batchSize)
	sucesses := 0
	failures := 0

	for i := 0; i < totalURLs; i += batchSize {
		end := i + batchSize
		if end > totalURLs {
			end = totalURLs
		}

		// Launch goroutines to fetch sources concurrently.
		for j := i; j < end; j++ {
			wg.Add(1)
			go source.FetchSource(sourceApiUrl, champions[j], &wg, result)
		}

		// Wait for all goroutines in this batch to finish.
		wg.Wait()

		// Collect results from channels
		for j := i; j < end; j++ {
			data := <-result
			if data != nil {
				err := dynamo.SaveProcessedCounters(dbInterface, tableName, data)
				if err != nil {
					log.Error().Err(err).Msg("failed to save " + data.Champion)
					failures++
				}
				sucesses++
			} else {
				failures++
			}
		}
		log.Info().Msgf("succeses: %d  failures: %d", sucesses, failures)
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
	refresh(champions.Champions)
	return nil
}

func offlineHandler() {
	test := []string{"ziggs"}
	refresh(test)
}
