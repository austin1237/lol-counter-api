package main

import (
	"context"
	"encoding/json"
	"fmt"
	"get/dynamo"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var tableName string

func init() {
	tableName = os.Getenv("COUNTER_TABLE_NAME")
	if tableName == "" {
		panic("Environment variable COUNTER_TABLE_NAME is not set or empty.")
	}
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" && os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		offlineHandler()
	} else {
		lambda.Start(handler)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	counter, err := dynamo.GetCounter(tableName, "ziggs")
	if err != nil {
		fmt.Println("Error retrieving item:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	if len(counter.Counters) == 0 {
		fmt.Println("Item not found.")
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	response, _ := json.Marshal(counter)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
	}, nil
}

func offlineHandler() {
	counter, err := dynamo.GetCounter(tableName, "ziggs")
	if err != nil {
		fmt.Println("Error retrieving item:", err)
		os.Exit(1)
	}

	if len(counter.Counters) == 0 {
		fmt.Println("Item not found.")
	} else {
		fmt.Println("Retrieved Item:", counter)
	}
}
