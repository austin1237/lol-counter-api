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

resource "aws_dynamodb_table" "lol_counters" {
  name           = "lol-counters-${terraform.workspace}"
  billing_mode   = "PAY_PER_REQUEST"  # On-demand capacity mode
  hash_key       = "Champion"

  attribute {
    name = "Champion"
    type = "S"  # String data type for Champion attribute
  }
}