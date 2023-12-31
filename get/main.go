package main

import (
	"context"
	"encoding/json"
	"get/dynamo"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog/log"
)

var tableName string
var dbInterface *dynamodb.DynamoDB

func init() {
	tableName = os.Getenv("COUNTER_TABLE_NAME")
	if tableName == "" {
		panic("Environment variable COUNTER_TABLE_NAME is not set or empty.")
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

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		offlineHandler()
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	champion := request.QueryStringParameters["champion"]
	counter, err := dynamo.GetCounter(dbInterface, tableName, champion)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving from dynamo item: " + champion)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	if len(counter.Counters) == 0 {
		message := "Champion: " + champion + " not found"
		log.Info().Msg(message)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       message,
		}, nil
	}

	response, _ := json.Marshal(counter)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
	}, nil
}

func offlineHandler() {
	champion := "ziggs"
	counter, err := dynamo.GetCounter(dbInterface, tableName, champion)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving from dynamo item: " + champion)
		os.Exit(1)
	}

	if len(counter.Counters) == 0 {
		log.Info().Msg("Champion: " + champion + " not found")
	} else {
		log.Info().Msgf("Retrieved Item:  %v", counter)
	}
}
