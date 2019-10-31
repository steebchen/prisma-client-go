FROM golang:1.13

WORKDIR /app

RUN curl -sL https://deb.nodesource.com/setup_10.x | bash - && apt-get install -y nodejs
RUN npm i -g prisma2 --unsafe-perm

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build .
RUN go generate ./...
RUN go test -v ./...
