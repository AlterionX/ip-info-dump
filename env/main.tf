terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

# Setting up basic information...

variable environment {
  type = string
}

variable geoipkey {
  type = string
}

variable virustotalkey {
  type = string
}

variable region {
  type = string
  default = "us-east-1"
}

provider "aws" {
  profile = "default"
  region = "${var.region}"
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

# Setting up identities...

# I could look into creating better roles if this project is a longer running one, but as it is a one off ...
resource "aws_iam_role" "lambda" {
  name = "ip_info_dump_iam"

  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  })
}

# Creating the lambda...
# I should hook this up with logging, but for now...

resource "aws_lambda_function" "ip_info_dump_lambda" {
  filename      = "../bin/main.zip"
  function_name = "ip_info_dump"
  role          = aws_iam_role.lambda.arn
  handler       = "main"

  source_code_hash = filebase64sha256("../bin/main.zip")

  runtime = "go1.x"

  environment {
    variables = {
      IPDUMP_GEOIP_KEY  = var.geoipkey
      VIRUSTOTAL_APIKEY = var.virustotalkey
    }
  }
}

# Setting up proxy...

resource "aws_api_gateway_rest_api" "gateway" {
  name        = "ip_info_dump_proxy"
  description = "Look up information about an IP quickly."
}

resource "aws_api_gateway_resource" "ip_info" {
  rest_api_id = "${aws_api_gateway_rest_api.gateway.id}"
  parent_id   = "${aws_api_gateway_rest_api.gateway.root_resource_id}"
  path_part   = "ip-info"
}

# I could probably set up an API key, but I'm still working on learning what usage plans are...
resource "aws_api_gateway_method" "get_ip_info_dump" {
  rest_api_id   = "${aws_api_gateway_rest_api.gateway.id}"
  resource_id   = "${aws_api_gateway_resource.ip_info.id}"
  http_method   = "GET"
  authorization = "NONE"
}

# And finally hook up the actual serverless function with the proxy.

resource "aws_lambda_permission" "get_ip_info_dump_auth" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ip_info_dump_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn    = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.gateway.id}/*/${aws_api_gateway_method.get_ip_info_dump.http_method}${aws_api_gateway_resource.ip_info.path}"
}

resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = aws_api_gateway_rest_api.gateway.id
  resource_id             = aws_api_gateway_resource.ip_info.id
  http_method             = aws_api_gateway_method.get_ip_info_dump.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.ip_info_dump_lambda.invoke_arn
}

