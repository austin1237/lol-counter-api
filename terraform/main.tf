terraform {
  backend "s3" {
    bucket         = "lol-counter-state"
    key            = "global/s3/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "lol-counter-state-lock"
    encrypt        = true
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.4"
    }
  }
  required_version = "~> 1.5"
}

# ---------------------------------------------------------------------------------------------------------------------
# Lambdas
# ---------------------------------------------------------------------------------------------------------------------

module "get_lambda" {
  source         = "./lambda"
  zip_location   = "../get/bootstrap.zip"
  name           = "get-lol-counter-${terraform.workspace}"
  handler        = "bootstrap"
  run_time       = "provided.al2"
  timeout        = 300
  dynamo_arn     = aws_dynamodb_table.lol_counters.arn
  env_vars = {
    "COUNTER_TABLE_NAME" = "${aws_dynamodb_table.lol_counters.name}"
  } 
}

module "ingest_lambda" {
  source         = "./lambda"
  zip_location   = "../ingest/bootstrap.zip"
  name           = "ingest-lol-counter-${terraform.workspace}"
  handler        = "bootstrap"
  run_time       = "provided.al2"
  timeout        = 300
  dynamo_arn     = aws_dynamodb_table.lol_counters.arn
  env_vars = {
    "COUNTER_TABLE_NAME" = "${aws_dynamodb_table.lol_counters.name}"
    "COUNTER_SOURCE_API_URL" = "${var.COUNTER_SOURCE_API_URL}"
    "BATCH_SIZE" = 30
  } 
}

# ---------------------------------------------------------------------------------------------------------------------
# Cloudwatch that will trigger the ingest lambda
# ---------------------------------------------------------------------------------------------------------------------
module "ingest_lambda_trigger" {
  source               = "./cloudwatch-lambda-trigger"
  # Every 3 days at 1pm MDT
  start_time           = "cron(0 19 */3 * ? *)"
  name                 = "ingest-lol-counter-trigger-${terraform.workspace}"
  lambda_function_name = "${module.ingest_lambda.name}"
  description          = "The timed trigger for ${module.ingest_lambda.name}"
  lambda_arn           = "${module.ingest_lambda.arn}"
}

# ---------------------------------------------------------------------------------------------------------------------
# API gateway
# ---------------------------------------------------------------------------------------------------------------------

module "counter_gateway" {
  source         = "./api-gateway"
  lambda_name =  "${module.get_lambda.name}"
  lambda_invoke_arn =   "${module.get_lambda.invoke_arn}"
  api_name = "lol-counter-${terraform.workspace}"
}

# ---------------------------------------------------------------------------------------------------------------------
#  DynamoDb
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_dynamodb_table" "lol_counters" {
  name           = "lol-counters-${terraform.workspace}"
  billing_mode   = "PAY_PER_REQUEST"  # On-demand capacity mode
  hash_key       = "Champion"

  attribute {
    name = "Champion"
    type = "S"  # String data type for Champion attribute
  }
}