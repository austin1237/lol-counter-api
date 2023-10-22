package dynamo

import (
	"ingest/source"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// SaveProcessedCounters saves the ProcessedCounters data to DynamoDB
func SaveProcessedCounters(tableName string, data *source.ProcessedCounters) error {
	// Create a new DynamoDB session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		// Endpoint: aws.String("http://localhost:8000"),
	})
	if err != nil {
		return err
	}

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	// Convert the ProcessedCounters struct to a DynamoDB attribute value map
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	// Add the current time
	av["lastUpdated"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	// Create the input parameters for the PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	// Save the data to DynamoDB
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
