FROM golang:1.13 as prepare

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build photongo
RUN go build -o /photongo .

COPY integration/ .
COPY . ./photongo

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/photongo/docker/lift.sh"]
