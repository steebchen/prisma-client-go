FROM golang:1.13 as build

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN cd test/integration/; go run github.com/prisma/prisma-client-go prefetch

RUN cd test/integration/; go run github.com/prisma/prisma-client-go db push --preview-feature --schema schemax.prisma

# generate the client in the integration folder
RUN cd test/integration/; go run github.com/prisma/prisma-client-go generate --schema schemax.prisma

# build the integration binary with all dependencies
RUN cd test/integration/; go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM golang:1.13

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/test/integration/dev.db /app/dev.db

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/main"]
