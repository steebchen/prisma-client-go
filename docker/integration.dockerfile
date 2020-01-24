FROM golang:1.13 as pre

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build prisma-client-go
RUN go build -o /prisma-client-go .

FROM golang:1.13 as build

WORKDIR /app

COPY --from=pre /prisma-client-go /prisma-client-go

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

COPY integration/ .
COPY . ./prisma-client-go

# generate the client in the integration folder
RUN go run github.com/prisma/prisma-client-go generate

# build the integration binary with all dependencies
RUN go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM ubuntu:16.04
# TODO try scratch image

WORKDIR /app

RUN apt-get update -qqy
RUN apt-get install -qqy openssl ca-certificates

COPY --from=build /app/main /app/main

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/main"]
