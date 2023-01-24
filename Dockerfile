FROM golang:1.19.3-alpine3.16

# Install build tools
RUN apk update \
        && apk add --no-cache --virtual build-tools curl

RUN apk add alpine-sdk

# Install air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go install gotest.tools/gotestsum@latest

# Purge build tools
RUN apk del --purge build-tools

WORKDIR /app