
FROM golang:1.16.2-buster

LABEL Maintainer="Kalai"

RUN apt-get update -y && \
    apt-get install zip -y

RUN go env
