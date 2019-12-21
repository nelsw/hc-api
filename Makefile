# This file was designed and developed to be used as an interface for making various aspects of this software library.
# For rapid development, use scripts to override environment variables when issuing make commands.
# In an effort to maintain modularity, fuctions separate concerns and knowledge by a domain driven data model.
# This not only supports AOP and separation of concerns, but organic atomiticy for clients and requests.

# The API domain function to make, see the README and ~go/src/hc-api/cmd/* for more information.
DOMAIN=
# The command, AKA switch statement case value, of an API request.
CMD=

# Build properties, effecitvely final.
SRC=main
SRC_EXE=${SRC}
SRC_ZIP=${SRC}.zip
SRC_DIR=./cmd/${DOMAIN}/${SRC}.go
ZIP_DIR=fileb://./${SRC_ZIP}
TST_OUT=cp.out
TST_DIR=./...

# SAM local invoke properties.
TMP_YML=template.yml
TMP_JSON=testdata/tmp.json
REQUEST_JSON=testdata/request.json
TEMPLATE_YML=testdata/template.yml
QSP=$(shell jq '.' testdata/${DOMAIN}/${CMD}/qsp.json -c)
BODY=$(shell jq '.|tostring' testdata/${DOMAIN}/${CMD}/body.json)

# AWS Lambda Function (λƒ) properties.
FUNCTION=handler
HANDLER=${SRC}
RUNTIME=go1.x
DESC="null"
TIMEOUT=30
MEMORY=512
ROLE=
ENVIRONMENT='$(shell jq '.' testdata/${DOMAIN}/env.json -c)'
VARIABLES=$(shell jq '.Variables' testdata/${DOMAIN}/env.json -c)

# A phony target is one that is not really the name of a file, but rather a sequence of commands.
# We use this practice to avoid potential naming conflicts with files in the home environment but
# also improve performance by telling the SHELL that we do not expect the command to create a file.
.PHONY: clean test build package invoke update create

# Convenience method for initializing request.json and templte.yml files prior to executing the invoke command.
it: init-request init-template invoke clean

# Removes build and package artifacts.
clean:
	rm -f ${SRC_ZIP}; rm -f ${SRC_EXE}; rm -f ${COVERAGE_REPORT}; rm -f ${TMP_YML}; rm -f ${TMP_JSON};

# Tests the entire project and outputs a coverage profile.
test:
	go test -coverprofile ${COVERAGE_REPORT} ${TST_DIR}

# Update the request event with test specific query string parameters and body data.
init-request:
	jq '.queryStringParameters=${QSP}' ${REQUEST_JSON} | sponge ${REQUEST_JSON};
	jq '.body=${BODY}' ${REQUEST_JSON} | sponge ${REQUEST_JSON};

# Update the sam template with domain and environment details.
init-template:
	jq -n '$(shell yq r -j ${TEMPLATE_YML})' > ${TMP_JSON};
	jq '.Resources.handler.Properties.Environment.Variables=${VARIABLES} | \
		.Resources.handler.Properties.Handler="${HANDLER}" | \
		.Resources.handler.Properties.MemorySize=${MEMORY} | \
		.Resources.handler.Properties.Runtime="${RUNTIME}" | \
		.Resources.handler.Properties.Timeout=${TIMEOUT} | \
		.Resources.handler.Properties.CodeUri="." | \
		.Description=${DESC}' ${TMP_JSON} | sponge ${TMP_JSON};
	yq r ${TMP_JSON} | sponge ${TMP_YML};

# Builds the source executable from a specified path.
build:
	GOOS=linux GOARCH=amd64 go build -o ${SRC_EXE} ${SRC_DIR}

# Packages the executable into a zip file with flags -9, compress better, and -r, recurse into directories.
package: build
	zip -9 -r ${SRC_ZIP} ${SRC_EXE}

# Executes build, package and `sam local invoke` with flags -t, path to required template.[yaml|yml] file, and -e,
# path to optional JSON file containing event data
# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-invoke.html
invoke: build package
	sam local invoke -t ${TMP_YML} -e ${REQUEST_JSON} ${FUNCTION} \
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
		--environment ${ENVIRONMENT} \
		--runtime ${RUNTIME}

# Helper command to update both code and configuration of our λƒ, and clean the projet directory.
update: update-code update-conf clean

# Creates an AWS λƒ.
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
		--environment ${ENVIRONMENT}
