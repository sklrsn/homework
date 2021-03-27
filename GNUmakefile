.DEFAULT_GOAL := all

NAME             = preprocessor
EXECUTEABLE_NAME = preprocessor
DESCRIPTION 	 = Parse SQS events and push it to Kinesis stream
DISTRIBUTION	 = linux
ARCH             = amd64

export LAMBDA_EXECUTOR = docker

.PHONY: all build package localstack env deploy clean

all: deps build package localstack env

deps:
	@echo "download dependencies ..."
	go mod vendor -v && \
	go mod tidy

build:
	@echo "compile binaries ..."
	@cd preprocessor && \
	env GOOS=${DISTRIBUTION} GOARCH=${ARCH} go build -ldflags="-s -w" -o dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME} .

package:
	@echo "package binaries ..."
	@cd preprocessor/dist/${DISTRIBUTION}/${ARCH} && \
	zip ${EXECUTEABLE_NAME}.zip ${EXECUTEABLE_NAME}

localstack:
	@echo "create Localstack environment ...$(LAMBDA_EXECUTOR)"
	@docker-compose up

env:
	aws --endpoint-url=http://localhost:4566 sqs list-queues
	aws --endpoint-url=http://localhost:4566 kinesis list-streams
	aws --endpoint-url=http://localhost:4566 lambda list-functions

deploy:
	@echo "deploy lambda ..."
	aws lambda create-function --function-name preprocessor --runtime go1.x \
	--zip-file fileb://preprocessor/dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}.zip \
	--handler preprocessor --endpoint-url=http://localhost:4566 \
	--role arn:aws:iam::skalai:role/execution_role

clean:
	@rm -rf vendor/
	@rm -rf preprocessor/dist
	@docker kill $$(docker ps -aq)
	@docker rm $$(docker ps -aq)
