# This file was designed and developed to be used as an interface for making various aspects of this software library.
# For rapid development, use scripts to override environment variables when issuing make commands.

# The API domain function to make, see the README and ~go/src/hc-api/cmd/* for more information.
DOMAIN=user
CMD=save

# Build properties, effecitvely final.
SRC=main
SRC_EXE=${SRC}
SRC_ZIP=${SRC}.zip
SRC_DIR=./cmd/${DOMAIN}/${SRC}.go
ZIP_DIR=fileb://./${SRC_ZIP}
TST_OUT=cp.out
TST_DIR=./...

# SAM local invoke properties.
REQUEST_JSON=request.json
TEMPLATE_YML=template.yml
QSP=$(shell jq '.' testdata/${DOMAIN}/${CMD}/qsp.json -c)
BODY=$(shell jq '.|tostring' testdata/${DOMAIN}/${CMD}/body.json)
TMP=testdata/tmp.json

# AWS Lambda Function (λƒ) properties.
FUNCTION=handler
HANDLER=${SRC}
RUNTIME=go1.x
DESC=Hemp Conduit API Handler for the ${DOMAIN} domain.
TIMEOUT=30
MEMORY=512
ROLE=arn:aws:iam::270722761968:role/service-role/test-role
ENV_VAR=$(shell jq '.Variables' testdata/env.json -c)

# A phony target is one that is not really the name of a file, but rather a sequence of commands.
# We use this practice to avoid potential naming conflicts with files in the home environment but
# also improve performance by telling the SHELL that we do not expect the command to create a file.
.PHONY: clean test build package invoke update create

# Convenience method for initializing request.json and templte.yml files prior to executing the invoke command.
it: init-request init-template invoke

# Removes build and package artifacts.
clean:
	rm -f ${TST_OUT}; rm -f ${SRC_ZIP}; rm -f ${SRC_EXE}; rm -f ${TMP}

# Tests the entire project and outputs a coverage profile.
test:
	go test -coverprofile ${TST_OUT} ${TST_DIR}

# Update the request event with test specific query string parameters and body data.
init-request:
	jq '.queryStringParameters=${QSP}' ${REQUEST_JSON} | sponge ${REQUEST_JSON};
	jq '.body=${BODY}' ${REQUEST_JSON} | sponge ${REQUEST_JSON};

# Update the sam template with domain and environment details.
init-template:
	jq -n '$(shell yq r -j template.yml)' > testdata/tmp.json;
	jq '.Resources.handler.Properties.Environment.Variables=${ENV_VAR} | \
		.Resources.handler.Properties.Handler="${HANDLER}" | \
		.Resources.handler.Properties.MemorySize=${MEMORY} | \
		.Resources.handler.Properties.Runtime="${RUNTIME}" | \
		.Resources.handler.Properties.Timeout=${TIMEOUT} | \
		.Resources.handler.Properties.CodeUri="." | \
		.Description="${DESC}"' ${TMP} | sponge ${TMP};
	yq r ${TMP} | sponge ${TEMPLATE_YML};

# Builds the source executable from a specified path.
build:
	GOOS=linux GOARCH=amd64 go build -o ${SRC_EXE} ${SRC_DIR}

# Packages the executable into a zip file.
# -9 compress better
# -r recurse into directories
package: build
	zip -9 -r ${SRC_ZIP} ${SRC_EXE}

# Executes test, build, package and `sam local invoke`.
# -t path to required template.[yaml|yml] file
# -e path to optional JSON file containing event data
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-invoke.html
invoke: build package
	sam local invoke -t ${TEMPLATE_YML} -e ${REQUEST_JSON} ${FUNCTION} \
	| jq '{statusCode: .statusCode, headers: .headers,  body: .body|fromjson}'

# Updates λƒ code with freshly packaged source.
# https://docs.aws.amazon.com/cli/latest/reference/lambda/update-function-code.html
update-code: package
	aws lambda update-function-code --function-name ${FUNCTION} --zip-file ${ZIP_DIR}

# Updates λƒ configuration with variable values.
# https://docs.aws.amazon.com/cli/latest/reference/lambda/update-function-configuration.html
update-conf:
	aws lambda update-function-configuration \
		--function-name ${FUNCTION} \
		--role ${ROLE} \
		--handler ${HANDLER} \
		--description ${DESC} \
		--timeout ${TIMEOUT} \
		--memory-size ${MEMORY} \
		--environment ${ENV_VAR} \
		--runtime ${RUNTIME}

# Helper command to update λƒ code and configuration.
update: update-code update-conf

# Creates an AWS Lambda Function (λƒ).
create: package
	aws lambda create-function \
		--function-name ${FUNCTION} \
		--runtime ${RUNTIME} \
		--role ${ROLE} \
		--handler ${HANDLER} \
		--description ${DESC} \
		--zip-file ${ZIP_DIR} \
		--memory-size ${MEMORY} \
		--timeout ${TIMEOUT} \
		--environment ${ENV_VAR}