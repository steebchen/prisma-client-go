FROM golang:1.13 as build

WORKDIR /app

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN cd integration/; go run github.com/prisma/prisma-client-go prefetch

RUN cd integration/; go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name init
RUN cd integration/; go run github.com/prisma/prisma-client-go migrate up --experimental

# generate the client in the integration folder
RUN cd integration/; go run github.com/prisma/prisma-client-go generate

# build the integration binary with all dependencies
RUN cd integration/; go build -o /app/main .

# start a new stage to test if the runtime fetching works
FROM golang:1.13
# TODO try scratch image. golang is used because it's available on both linux and windows

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/integration/dev.db /app/dev.db

ENV PHOTON_GO_LOG=info
ENV DEBUG=*

CMD ["/app/main"]
