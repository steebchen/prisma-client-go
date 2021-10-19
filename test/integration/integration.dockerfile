FROM golang:1.16 as build

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

COPY . ./

WORKDIR /app/test/integration

RUN go get github.com/prisma/prisma-client-go@main

RUN go run github.com/prisma/prisma-client-go migrate dev --name init

RUN go mod tidy

RUN pwd
RUN ls -l
RUN cat go.mod
RUN cat go.sum

# build the integration binary with all dependencies
RUN go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM golang:1.16

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/test/integration/dev.db /app/dev.db

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/main"]
