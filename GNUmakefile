.DEFAULT_GOAL := all

NAME             = preprocessor
EXECUTEABLE_NAME = preprocessor
DESCRIPTION 	 = Parse SQS events and push it to Kinesis stream
DISTRIBUTION	 = linux
ARCH             = amd64

.PHONY: all build package localstack

all: deps build package localstack

deps:
	@echo "download dependencies ..."
	go mod vendor -v && \
	go mod tidy

build:
	@echo "compile binaries ..."
	@cd preprocessor && \
	GOOS=${DISTRIBUTION} GOARCH=${ARCH} go build -o dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME} .

package:
	@echo "package binaries ..."
	@cd preprocessor && \
	zip dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}.zip dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}

localstack:
	@echo "create Localstack environment ..."
	@docker-compose up -d

clean:
	@rm -rf vendor/
	@rm -rf preprocessor/dist
	@docker kill $$(docker ps -aq)
	@docker rm $$(docker ps -aq)
