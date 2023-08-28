# AWS go lambdas running on provided.al2 runtime have to be called bootstrap

packageLambdas:
	cd get && GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go && zip bootstrap.zip bootstrap && rm bootstrap && ls
