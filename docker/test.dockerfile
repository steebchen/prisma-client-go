FROM golang:1.13

WORKDIR /app

RUN curl -sL https://deb.nodesource.com/setup_10.x | bash - && apt-get install -y nodejs

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# make sure to always use the latest prisma2 to detect breaking changes for now
RUN npm i -g prisma2@alpha --unsafe-perm

RUN go build .
RUN go generate ./...
RUN go test -v ./...
