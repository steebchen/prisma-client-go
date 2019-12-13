FROM golang:1.13

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

RUN go build .
RUN go generate ./...
RUN go test -v ./...
