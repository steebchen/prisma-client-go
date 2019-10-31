FROM golang:1.13

USER root

WORKDIR /app

RUN go get -u golang.org/x/lint/golint

COPY . ./

RUN go fmt ./...
RUN golint ./...
