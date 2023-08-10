variable "zip_location" {
  description = "path to the ziped lambda"
}

variable "name" {
  description = "The name of the lambda function"
}

variable "handler" {
  description = "name of the lambdas handler"
}

variable "run_time" {
  description = "run time of the lambda"
}

variable "env_vars" {
  type        = map(string)
  description = "run time of the lambda"
}

variable "memory_size" {
  description = "Amount of memory in MB your Lambda Function can use at runtime. CPU is implicitly tied to this."
  default     = 128
}

variable "timeout" {
  description = "The max number of seconds the lambda can run"
  default     = 3
}

variable "dynamo_arn" {
  description = "arn of the dynamo table that this lambda will be given access to"
}