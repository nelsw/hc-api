.PHONY: clean build test package invoke update create call

FUNCTION_NAME=
EVENT=find

clean:
	rm -f main;
	rm -f main.zip;

test:
	USER_TABLE=${USER_TABLE}
	go test -coverprofile cp.out ./...

build:
	GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

package: build
	zip -9 -r main.zip main

invoke: build
	sam local invoke -e testdata/${EVENT}.json ${NAME} | jq

update: package
	aws lambda update-function-code --function-name ${NAME} --zip-file fileb://./handler.zip;\
	aws lambda update-function-configuration --function-name ${NAME} --handler main

create: package
	aws lambda create-function \
		--function-name ${NAME} \
		--description \
		--role ${ROLE} \
		--zip-file fileb://./handler.zip \
		--handler handler \
		--runtime go1.x \
		--memory-size 512 \
		--timeout 30;

call:
	curl -v "https://i5a2n7eqb0.execute-api.us-east-1.amazonaws.com/dev?cmd=find" | jq;