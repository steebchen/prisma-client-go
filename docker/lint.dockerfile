FROM golang:1.13 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build .
RUN go generate ./...

FROM golangci/golangci-lint:v1.21.0

WORKDIR /app

COPY --from=build /app /app

RUN golangci-lint run ./... -v --enable "gofmt,golint,scopelint,gocritic"
