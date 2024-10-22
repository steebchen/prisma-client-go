# Deploying Best Practices

## Best practices

Prisma Client Go needs to generate code to work. This means that you need to run the generator before you can use it. We
recommend to use Go's built-in [`go generate`](https://blog.golang.org/generate) to run the generator.

### Set up go generate

While this step is not required, we recommend to use [`go generate`](https://blog.golang.org/generate) to simplify
generating the Prisma Client Go client. To do so, you can just put the following line into a go file, and then run go
generate so `go run github.com/steebchen/prisma-client-go generate` and any other generate commands you run will get
executed.

Put this line into a Go file in your project, usually in `main.go`:

```go
//go:generate go run github.com/steebchen/prisma-client-go generate

func main() {
// ...
}
```

Now, run `go generate`:

```shell script
go generate ./...
```

Your Prisma Client Go code is now generated.
