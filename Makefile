# AWS go lambdas running on provided.al2 runtime have to be called bootstrap 
# https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-naming

packageLambdas: packageGet packageIngest
	
packageGet:
	cd get && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls

packageIngest:
	cd ingest && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls

