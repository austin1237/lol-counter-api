resource "aws_iam_policy" "lambda_logs_policy" {
  name        = "lambda-cloudwatch-logs-policy-${var.name}"
  description = "Allows Lambda function to write logs to CloudWatch Logs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda_dynamodb_policy" {
  name = "LambdaDynamoDBPolicy-${terraform.workspace}-${var.name}"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "dynamodb:GetItem",
          "dynamodb:Query",
          "dynamodb:Scan",
        ],
        Resource = [
          "${var.dynamo_arn}",
        ],
      },
    ],
  })
}

resource "aws_iam_role" "lambda_role" {
  name = "my-lambda-role-${terraform.workspace}-${var.name}"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}


resource "aws_iam_role_policy_attachment" "lambda_dynamodb_policy_attachment" {
  policy_arn = aws_iam_policy.lambda_dynamodb_policy.arn
  role       = aws_iam_role.lambda_role.name
}

resource "aws_iam_role_policy_attachment" "lambda_logs_attachment" {
  policy_arn = aws_iam_policy.lambda_logs_policy.arn
  role       = aws_iam_role.lambda_role.name
}

resource "aws_lambda_function" "lambda" {
  filename      = "${var.zip_location}"
  function_name = "${var.name}"
  role          = "${aws_iam_role.lambda_role.arn}"
  handler       = "${var.handler}"
  timeout       = "${var.timeout}"
  memory_size   = "${var.memory_size}"

  source_code_hash = "${filesha256("${var.zip_location}")}"
  runtime          = "${var.run_time}"

  environment {
    variables = "${var.env_vars}"
  }
}