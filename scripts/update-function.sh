#!/usr/bin/env bash
#
# This script validates an AWS Lambda Function code and configuration by providing valid variable parameters.
echo "==> Updating an AWS Lambda Function..."

# The following conditions validation command variable by asserting and defaulting values.
if [ -z "${DOMAIN}" ]; then echo "ERROR: Set DOMAIN to the entity identity of the λƒ to build"; exit 1; fi
if [ -z "${FUNCTION}" ]; then echo "ERROR: Set FUNCTION to the name of the λƒ to build"; exit 1; fi
if [ -z "${ROLE}" ]; then echo "ERROR: Set ROLE to the AWS IAM Role of the λƒ to build"; exit 1; fi
if [ -z "${TIMEOUT}" ]; then TIMEOUT="30"; fi
if [ -z "${MEMORY}" ]; then MEMORY="512"; fi
if [ -z "${DESC}" ]; then DESC="null"; fi

# https://docs.aws.amazon.com/cli/latest/reference/lambda/update-function-configuration.html
aws lambda update-function-configuration \
  --function-name "${FUNCTION}" \
  --role "${ROLE}" \
  --description "${DESC}" \
  --timeout "${TIMEOUT}" \
  --memory-size "${MEMORY}" \
  --environment "$(jq '.' test/"${DOMAIN}"/env.json -c)";

# https://docs.aws.amazon.com/cli/latest/reference/lambda/update-function-code.html
aws lambda update-function-code \
  --function-name "${FUNCTION}" \
  --zip-file fileb://./main.zip

echo "==> Function Updated!"
