.DEFAULT_GOAL := all

NAME             = preprocessor
EXECUTEABLE_NAME = preprocessor
DESCRIPTION 	 = Parse SQS events and push it to Kinesis stream
DISTRIBUTION	 = linux
ARCH             = amd64

export LAMBDA_EXECUTOR = docker
export TMPDIR=/tmp

.PHONY: all compile run env clean

all: run env

deps:
	@echo "download go dependencies ..."
	@cd preprocessor && \
	go mod vendor -v && \
	go mod tidy

compile:
	@echo "compile binaries ..."
	@cd preprocessor && \
	GOOS=${DISTRIBUTION} GOARCH=${ARCH} go build -ldflags="-s -w" -o dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME} .

package:
	@echo "package binaries ..."
	@cd preprocessor/dist/${DISTRIBUTION}/${ARCH} && \
	zip ${EXECUTEABLE_NAME}.zip ${EXECUTEABLE_NAME}

run:
	@echo "create localstack environment and launch application ..."
	@echo LAMBDA_EXECUTOR ...$(LAMBDA_EXECUTOR)
	@echo TMPDIR ...$(TMPDIR)
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
	@rm -rf preprocessor/vendor/
	@rm -rf preprocessor/dist
	@docker rmi -f homework_preprocessor
	@docker rmi -f homework_localstack
	@docker rmi -f homework_sensor-fleet
