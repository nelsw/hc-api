#!/usr/bin/env bash
#
# This script creates an AWS Lambda Function.
echo "==> Creating an AWS Lambda Function..."

if [ -z "${FUNCTION}" ]; then echo "ERROR: Set FUNCTION to the function name of the λƒ to build"; exit 1; fi
if [ -z "${DOMAIN}" ]; then echo "ERROR: Set DOMAIN to the entity identity of the λƒ to build"; exit 1; fi
if [ -z "${ROLE}" ]; then echo "ERROR: Set ROLE to the AWS IAM Role of the λƒ to build"; exit 1; fi
if [ -z "${TIMEOUT}" ]; then TIMEOUT="30"; fi
if [ -z "${MEMORY}" ]; then MEMORY="512"; fi
if [ -z "${DESC}" ]; then DESC="null"; fi

# https://docs.aws.amazon.com/cli/latest/reference/lambda/create-function.html
aws lambda create-function \
  --function-name "${FUNCTION}" \
  --runtime "go1.x" \
  --role "${ROLE}" \
  --handler "main" \
  --description "${DESC}" \
  --zip-file "fileb://./main.zip" \
  --memory-size "${MEMORY}" \
  --timeout "${TIMEOUT}" \
  --environment "$(jq '.' test/"${DOMAIN}"/env.json -c)"

echo "==> Function Created!"