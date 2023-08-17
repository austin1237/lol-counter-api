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