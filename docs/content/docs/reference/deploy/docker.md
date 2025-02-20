# Deploy Via Docker

## Docker

To deploy with docker, build your dockerfile as usual, run `go generate ./...` (
see [setting up go generate](best-practices#set-up-go-generate)), and you're good to go!

We also recommend using [Go modules](https://blog.golang.org/using-go-modules), which is recommended when using Go >
=1.13.

## Example Dockerfile

Your dockerfile could look like this. It uses Go modules and layered caching for fast docker builds.

If you want to optimize your docker images even further, learn how to use a [multi-stage build](#optimized-dockerfile).

```dockerfile
FROM golang:1.20 as build

WORKDIR /workspace

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch

COPY . ./

# generate the Prisma Client Go client
RUN go run github.com/steebchen/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /app .

CMD ["/app"]
```

## Optimized Dockerfile

If you want to optimize your docker images even further, you can use a
multi-stage build. This will create a smaller image, which only contains the
binary and not the whole build environment. However, you need to have some
extra steps in place, such as copying SSL certificates.

```dockerfile
FROM golang:1.21.5-buster as builder

WORKDIR /workspace

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch

COPY ./ ./
# generate the Prisma Client Go client
RUN go run github.com/steebchen/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

# use the scratch image for the smallest possible image size
FROM scratch

# copy over SSL certificates, so that we can make HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /workspace/app /app

ENTRYPOINT ["/app"]

```
