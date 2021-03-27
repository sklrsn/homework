.DEFAULT_GOAL := all

NAME             = preprocessor
EXECUTEABLE_NAME = preprocessor
DESCRIPTION 	 = Parse SQS events and push it to Kinesis stream
DISTRIBUTION	 = linux
ARCH             = amd64


all: build package

build:
	cd preprocessor && \
	GOOS=${DISTRIBUTION} GOARCH=${ARCH} go build -o dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME} .

package:
	cd preprocessor && \
	zip dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}.zip dist/${DISTRIBUTION}/${ARCH}/${EXECUTEABLE_NAME}

clean:
	rm -rf preprocessor/dist