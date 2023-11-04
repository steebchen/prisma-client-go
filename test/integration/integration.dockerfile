ARG IMAGE
FROM $IMAGE as build

WORKDIR /app

RUN go version

ENV PRISMA_CLIENT_GO_LOG=info
ENV DEBUG=*

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./

WORKDIR /app/test/integration

RUN go run github.com/steebchen/prisma-client-go db push

# build the integration binary with all dependencies
RUN go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM $IMAGE

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/test/integration/dev.db /app/dev.db

ENV PRISMA_CLIENT_GO_LOG=debug
ENV DEBUG=*

CMD ["/app/main"]
