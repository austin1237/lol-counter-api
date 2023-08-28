# lol-counter-api
An api/script that consumes leauge and of legends counter data from [lol-counter-source-api](https://github.com/austin1237/lol-counter-source-api) and exposes for faster consumption.

## Prerequisites
You must have the following installed/configured on your system for this to work correctly<br />
1. [Go](https://go.dev/doc/install)

## Environment Variables
The following variables need to be set on your local/ci system.

### SOURCE_API_URL
Url of the deployed lol-counter-source-api api
### COUNTER_TABLE_NAME
name of the dynamodb table that will be used

## Deployment
Deployment currently uses [Terraform](https://www.terraform.io/) to set up AWS services.

### Setting up remote state
Terraform has a feature called [remote state](https://www.terraform.io/docs/state/remote.html) which ensures the state of your infrastructure to be in sync for mutiple team members as well as any CI system.

This project **requires** this feature to be configured. To configure **USE THE FOLLOWING COMMAND ONCE PER TEAM**.

```bash
cd terraform/remote-state
terraform init
terraform apply
```


