FROM golang:1.20.1 as build

WORKDIR /app

ENV PRISMA_CLIENT_GO_LOG=info
ENV DEBUG=*

COPY . ./

WORKDIR /app/test/integration

RUN go mod download -x

RUN go run github.com/prisma/prisma-client-go db push --schema schemax.prisma

# build the integration binary with all dependencies
RUN go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM golang:1.20.1

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/test/integration/dev.db /app/dev.db

ENV PRISMA_CLIENT_GO_LOG=info
ENV DEBUG=*

CMD ["/app/main"]
