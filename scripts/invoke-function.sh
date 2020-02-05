#!/usr/bin/env bash
#
# This script simulates an AWS Lambda Function invocation using the serverless application model framework.
echo "==> Invoking a local AWS Lambda Function..."

# The following conditions validation required command variables.
if [ -z "${FN}" ]; then echo "ERROR: Set FN to the function name of the λƒ to build"; exit 1; fi
if [ -z "${ROLE}" ]; then echo "ERROR: Set ROLE to the AWS IAM Role of the λƒ to build"; exit 1; fi
if [ -z "${DOM}" ]; then echo "ERROR: Set DOM to the domain of the λƒ to build"; exit 1; fi
if [ -z "${CASE}" ]; then echo "ERROR: Set CASE to the action performed by the λƒ to build"; exit 1; fi

# Create a temporary template.json file for local sam invocation.
jq -n "$(yq r -j build/template.yml)" > template.json;

# Update the sam template with required properties.
if [ -z "${TIMEOUT}" ];  then jq --arg var "${TIMEOUT}" '.Resources.handler.Properties.Timeout=$var' template.json; fi
if [ -z "${MEMORY}" ]; then jq --arg var "${MEMORY}" '.Resources.handler.Properties.MemorySize=$var' template.json; fi
if [ -z "${DESC}" ]; then jq --arg var "${DESC}" '.Description=$var' template.json; fi

# Update the sam template with environment variables.
jq --arg var "$(shell jq '.Variables' test/"${DOMAIN}"/env.json -c)" \
'.Resources.handler.Properties.Environment.Variables=$var' template.json;

# Create a temporary request.json file for local sam invocation.
jq -n "$(yq r -j build/request.json)" > request.json;
jq --arg var "$(shell jq '.|tostring' test/"${DOMAIN}"/"${CMD}"/body.json)" '.body=$var' request.json

# Execute `sam local invoke` with flags required template file, and optional file containing event data.
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-invoke.html
sam local invoke -t template.json -e request.json "${FUNCTION}" | \
jq '{statusCode: .statusCode, headers: .headers,  body: .body|fromjson}'

echo "==> Function Invoked!"