FROM golang:1.13 as prepare

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/docker/migrate.sh"]
