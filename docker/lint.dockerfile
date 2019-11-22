FROM golang:1.13 as build

WORKDIR /app

RUN curl -sL https://deb.nodesource.com/setup_10.x | bash - && apt-get install -y nodejs
RUN npm i -g prisma2@alpha --unsafe-perm

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build .
RUN go generate ./...

FROM golangci/golangci-lint:v1.21.0

WORKDIR /app

COPY --from=build /app /app

RUN golangci-lint run ./... -v --enable "gofmt,golint,scopelint,gocritic"
