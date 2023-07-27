#!/bin/bash

# Wait for DynamoDB Local to be ready
# until aws dynamodb list-tables --endpoint-url http://localhost:8000 2>/dev/null; do
#     sleep 1
# done

# Create the "lol-counters" table
aws dynamodb create-table --table-name lol-counters \
    --attribute-definitions AttributeName=Champion,AttributeType=S \
    --key-schema AttributeName=Champion,KeyType=S \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:8000
    --output json

# Wait for the table to be active
aws dynamodb wait table-exists --table-name lol-counters --endpoint-url http://localhost:8000 --output json

echo "command finsihed"
