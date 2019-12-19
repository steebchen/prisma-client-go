FROM golang:1.13 as build

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod ./
RUN go mod download

COPY . ./

# build photongo
RUN go build .

# generate photon in integration folder
RUN cd integration/ && go run github.com/prisma/photongo generate

# build the integration binary with all dependencies
RUN cd integration/ && go build -o /main .

# start a new stage to test if the runtime fetching works
FROM ubuntu:16.04

RUN apt-get update -qqy
RUN apt-get install -qqy openssl ca-certificates

COPY --from=build /main /main

CMD ["/main"]
