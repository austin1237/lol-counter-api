# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

variable "COUNTER_SOURCE_API_URL" {
    sensitive = true
}

variable "INGEST_CRON" {
}