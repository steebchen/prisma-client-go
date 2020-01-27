FROM golang:1.13 as prepare

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build prisma-client-go
RUN go build -o /prisma-client-go .

COPY integration/ .
COPY . ./prisma-client-go

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/prisma-client-go/docker/lift.sh"]
