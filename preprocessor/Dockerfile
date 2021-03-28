
FROM golang:latest as builder

LABEL Maintainer="Kalai"

RUN go env

WORKDIR /go/src/github.com/sklrsn/homework/preprocessor

COPY . .

RUN go mod vendor -v && \
    go mod tidy && \
    go build -o dist/linux/amd64/preprocessor .

FROM builder as app

COPY --from=builder \
    /go/src/github.com/sklrsn/homework/preprocessor/dist/linux/amd64/preprocessor /opt/bin/preprocessor

RUN chmod a+x /opt/bin/preprocessor

CMD [ "/opt/bin/preprocessor" ]
