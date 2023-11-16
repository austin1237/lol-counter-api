package dynamo

import (
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDB is a mock implementation for DynamoDBAPI
type MockDynamoDB struct {
	mock.Mock
	DynamoDBAPI
}

// GetItem is a mocked implementation for the GetItem method
func (m *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func TestGetCounter(t *testing.T) {
	// Create a mock instance
	mockDynamoDB := new(MockDynamoDB)

	// Set expectations on the mocked GetItem method
	mockDynamoDB.On("GetItem", mock.Anything).
		Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"Champion":    {S: aws.String("MockChampion")},
				"lastUpdated": {N: aws.String(strconv.FormatInt(time.Now().Unix(), 10))},
				"Counters":    {L: []*dynamodb.AttributeValue{{S: aws.String("Counter1")}}},
			},
		}, nil)

	_, err := GetCounter(mockDynamoDB, "mockTable", "mockChampion")

	// Call the method under test

	// Assert the result or any other expectations
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the expected method was called
	mockDynamoDB.AssertExpectations(t)
}
