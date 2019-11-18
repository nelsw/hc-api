# This file was designed and developed to be used as an interface for making various aspects of this software library.
# For rapid development, use scripts to override environment variables when issuing make commands.

# The API domain to make, see the the subdirecgories of cmd for more informaiton.
DOMAIN=

# Build properties, effecitvely final.
SRC=main
SRC_EXE=${SRC}
SRC_ZIP=${SRC}.zip
SRC_DIR=./cmd/${DOMAIN}/${SRC}.go
ZIP_DIR=fileb://./${SRC_ZIP}
TST_OUT=cp.out
TST_DIR=./...
TEMPLATE_YML_DIR=template.yml
REQUEST_JSON_DIR=request.json

# AWS Lambda Function (λƒ) properties.
FUNCTION=
HANDLER=
RUNTIME=
DESC=
TIMEOUT=
MEMORY=
ROLE=
ENV_VAR=

# A phony target is one that is not really the name of a file, but rather a sequence of commands.
# We use this practice to avoid potential naming conflicts with files in the home environment but
# also improve performance by telling the SHELL that we do not expect the command to create a file.
.PHONY: clean test build package invoke update create

# Removes build and package artifacts.
clean:
	rm -f ${TST_OUT};
	rm -f ${SRC_ZIP};
	rm -f ${SRC_EXE};

# Tests the entire project and outputs a coverage profile.
test:
	go test -coverprofile ${TST_OUT} ${TST_DIR}

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
invoke: test build package
	sam local invoke \
		-t ${TEMPLATE_YML_DIR} \
		-e ${REQUEST_JSON_DIR} \
		${FUNCTION}

# Updates λƒ code with freshly packaged source.
update-code: package
	aws lambda update-function-code \
		--function-name ${FUNCTION} \
		--zip-file ${ZIP_DIR}

# Updates λƒ  configuration with variable values.
update-conf:
	aws lambda update-function-configuration \
		--function-name ${FUNCTION} \
		--role ${ROLE} \
		--handler ${HANDLER} \
		--description ${DESC}
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
		--role ${ROLE}
		--handler ${HANDLER} \
		--description ${DESC} \
		--zip-file ${ZIP_DIR} \
		--memory-size ${MEMORY} \
		--timeout ${TIMEOUT} \
		--environment ${ENV_VAR}