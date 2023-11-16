package dynamo

import (
	"ingest/source"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDB is a mock implementation for DynamoDBAPI
type MockDynamoDB struct {
	mock.Mock
	DynamoDBAPI
}

// PutItem is a mocked implementation for the PutItem method
func (m *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func TestSaveProcessedCounters(t *testing.T) {
	// Create a mock instance
	mockDynamoDB := new(MockDynamoDB)

	// Set expectations on the mocked PutItem method
	mockDynamoDB.On("PutItem", mock.Anything).
		Return(&dynamodb.PutItemOutput{}, nil)

	mockData := &source.ProcessedCounters{}

	err := SaveProcessedCounters(mockDynamoDB, "mock_tableName", mockData)

	// Assert the result or any other expectations
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert that the expected method was called
	mockDynamoDB.AssertExpectations(t)
}
