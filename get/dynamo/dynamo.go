package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ProcessedCounters struct {
	Champion string   `json:"champion"`
	Counters []string `json:"counters"`
}

func GetCounter(tablename string, champion string) (ProcessedCounters, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		// Endpoint: aws.String("http://localhost:8000"),
	})

	if err != nil {
		return ProcessedCounters{}, err
	}

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"Champion": {
				S: aws.String(champion), // Replace "YourPrimaryKeyAttributeName" with your actual primary key attribute name
			},
		},
	}

	// Get the item from DynamoDB
	result, err := svc.GetItem(input)

	if err != nil {
		return ProcessedCounters{}, err
	}

	if result.Item == nil {
		return ProcessedCounters{}, nil
	}

	// Convert the DynamoDB item to ProcessedCounters struct
	counters := ProcessedCounters{
		Champion: *result.Item["Champion"].S,
	}

	// Unpack the Counters attribute (which is a DynamoDB List) into a slice of strings
	if result.Item["Counters"] != nil && result.Item["Counters"].L != nil {
		for _, av := range result.Item["Counters"].L {
			counters.Counters = append(counters.Counters, *av.S)
		}
	}

	return counters, nil
}