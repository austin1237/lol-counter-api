resource "aws_api_gateway_rest_api" "counter_source" {
  name = "${var.api_name}"
}

resource "aws_api_gateway_resource" "counter_path" {
  rest_api_id = aws_api_gateway_rest_api.counter_source.id
  parent_id   = aws_api_gateway_rest_api.counter_source.root_resource_id
  path_part   = "counter"
}

resource "aws_api_gateway_method" "counter_method" {
  rest_api_id   = aws_api_gateway_rest_api.counter_source.id
  resource_id   = aws_api_gateway_resource.counter_path.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "counter_integration" {
  rest_api_id             = aws_api_gateway_rest_api.counter_source.id
  resource_id             = aws_api_gateway_resource.counter_path.id
  http_method             = aws_api_gateway_method.counter_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${var.lambda_invoke_arn}"
}

resource "aws_api_gateway_deployment" "counter_deployment" {
  depends_on = [aws_api_gateway_integration.counter_integration]
  rest_api_id = aws_api_gateway_rest_api.counter_source.id
  stage_name  = "prod"
  triggers = {

    # NOTE: The configuration below will satisfy ordering considerations,
    #       but not pick up all future REST API changes. More advanced patterns
    #       are possible, such as using the filesha1() function against the
    #       Terraform configuration file(s) or removing the .id references to
    #       calculate a hash against whole resources. Be aware that using whole
    #       resources will show a difference after the initial implementation.
    #       It will stabilize to only change when resources change afterwards.
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.counter_path.id,
      aws_api_gateway_method.counter_method.id,
      aws_api_gateway_integration.counter_integration.id,
    ]))
  }
}

resource "aws_api_gateway_method_settings" "general_settings" {
  rest_api_id = "${aws_api_gateway_rest_api.counter_source.id}"
  stage_name  = "${aws_api_gateway_deployment.counter_deployment.stage_name}"
  method_path = "*/*"

  settings {
    # Enable CloudWatch logging and metrics
    metrics_enabled        = true
    data_trace_enabled     = true
    logging_level          = "INFO"

    # Limit the rate of calls to prevent abuse and unwanted charges
    throttling_rate_limit  = 100
    throttling_burst_limit = 50
  }
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "AllowMyAPIInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${var.lambda_name}"
  principal     = "apigateway.amazonaws.com"

  # The /* part allows invocation from any stage, method and resource path
  # within API Gateway.
  source_arn = "${aws_api_gateway_rest_api.counter_source.execution_arn}/*"
}