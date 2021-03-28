.DEFAULT_GOAL := all

NAME             = preprocessor
EXECUTEABLE_NAME = preprocessor
DESCRIPTION 	 = Parse SQS events and push it to Kinesis stream
DISTRIBUTION	 = linux
ARCH             = amd64

export LAMBDA_EXECUTOR = docker
export TMPDIR=/tmp

BUILDER_IMAGE	= builder:latest

.PHONY: all compile run env clean builder lambda serverless build

all: run env

deps:
	@echo "download go dependencies ..."
	@cd preprocessor && \
	go mod vendor -v && \
	go mod tidy

builder:
	docker build -t $(BUILDER_IMAGE) .

build:
	docker run -it --rm --net host \
		-e DISTRIBUTION=${DISTRIBUTION} \
		-e ARCH=${ARCH} \
		-e DISTRIBUTION=${DISTRIBUTION} \
		-e EXECUTEABLE_NAME=${EXECUTEABLE_NAME} \
		-v ${PWD}:/homework \
		-w /homework \
		--entrypoint /bin/bash $(BUILDER_IMAGE) \
		-c 'make compile; make package'
compile:
	@echo "compile binaries ..."
	@cd preprocessor && \
	GOOS=${DISTRIBUTION} GOARCH=${ARCH} go build -ldflags="-s -w" -o dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME} .
	@echo "binaries written to preprocessor/dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}"

package:
	@echo "package binaries ..."
	@cd preprocessor/dist/${DISTRIBUTION}/${ARCH} && \
	zip ${EXECUTEABLE_NAME}.zip ${EXECUTEABLE_NAME}
	@echo "preprocessor.zip available at preprocessor/dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}"

run:
	@echo "create localstack environment and launch application ..."
	@echo LAMBDA_EXECUTOR ...$(LAMBDA_EXECUTOR)
	@echo TMPDIR ...$(TMPDIR)
	@docker-compose up --build

env:
	aws --endpoint-url=http://localhost:4566 sqs list-queues
	aws --endpoint-url=http://localhost:4566 kinesis list-streams
	aws --endpoint-url=http://localhost:4566 lambda list-functions

lambda:
	@echo "deploying lambda to localstack ..."
	aws lambda create-function --function-name preprocessor --runtime go1.x \
	--zip-file fileb://preprocessor/dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}.zip \
	--handler preprocessor --endpoint-url=http://localhost:4566 \
	--role arn:aws:iam::kalai:role/execution_role

	@echo "enable trigger to launch lambdas when message published to SQS..."
	aws --endpoint-url=http://localhost:4566 lambda create-event-source-mapping \
	--event-source-arn arn:aws:sqs:eu-west-1:000000000000:submissions \
	--function-name preprocessor

serverless: builder build env lambda env

clean:
	@rm -rf preprocessor/vendor/
	@docker rmi -f homework_preprocessor
	@docker rmi -f homework_localstack
	@docker rmi -f homework_sensor-fleet
