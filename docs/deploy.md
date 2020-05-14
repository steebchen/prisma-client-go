
### Deploy

Deploying a Prisma Client Go adds a few more steps, because it depends on the Prisma query engine, which is a binary we automatically download in your project path. Depending on where you deploy your code to, you might need to follow some extra steps.

#### Set up go generate

While this step is not required, we recommend to use [`go generate`](https://blog.golang.org/generate) to simplify generating the Prisma Client Go client. To do so, you can just put the following line into a go file, and then run go generate so `go run github.com/prisma/prisma-client-go generate` and any other generate commands you run will get executed.

Put this line into a go file in your project, usually in `main.go`:

```go
//go:generate go run github.com/prisma/prisma-client-go generate

func main() {
  // ...
}
```

Now, run `go generate`:

```shell script
go generate ./...
```

Your Prisma Client Go code is now generated.

#### Traditionally deploy to a server

Usually, you would deploy your Go app by running `go build .`, which generates a binary, and then deploy that binary anywhere you want. However, since Prisma Client Go depends on the Prisma query engine, you also need to deploy the query engine binary `query-engine-*` files.

This engine gets fetched at runtime per default for getting started quickly. However, in production deployments, you may want to prefetch this binary to have reproducible builds and/or images.

If you use different development environments, e.g. a Mac to develop, and Debian on your server, you need to specify these two binaries in the schema.prisma file so that you can then also upload the binary suitable for your deploy environment.

```prisma
generator db {
  provider = "go run github.com/prisma/prisma-client-go"
  binaryTargets = ["native", "debian-openssl-1.1.x"]
}
```

You can find all binary targets [in our specs repository](https://github.com/prisma/specs/tree/master/binaries#binary-builds).

#### Using docker

When deploying with docker, the setup is super easy. Build your dockerfile as usual, run `go generate ./...` (see [setting up go generate](#set-up-go-generate)), and you're good to go!

We also recommend using [Go modules](https://blog.golang.org/using-go-modules), which is recommended when using Go >=1.13.

Your dockerfile could look like this. It uses Go modules, layered caching for fast docker builds and multiple stages for lightweight images (usually a few megabytes).

```dockerfile
FROM golang:1.13 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

COPY . ./

# generate the Prisma Client Go client
RUN go run github.com/prisma/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /main .

CMD ["/main"]
```
