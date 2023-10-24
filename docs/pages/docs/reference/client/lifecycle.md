# Client lifecycle

## Connecting

To create a new Prisma client instance:

```go
client := db.NewClient()
if err := client.Prisma.Connect(); err != nil {
  handle(err)
}
```

## Disconnecting

Ideally, you should disconnect from the database when you're done:

```go
if err := client.Prisma.Disconnect(); err != nil {
  panic(fmt.Errorf("could not disconnect: %w", err))
}
```

If you're using a webserver, the best to handle it is to catch the `SIGTERM` signal, disconnect from the database and
afterwards clean up your webserver:

```go
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
go func() {
  <-c
  if err := client.Prisma.Disconnect(); err != nil {
    panic(fmt.Errorf("could not disconnect: %w", err))
  }
  // clean up your webserver here
  // e.g. httpServer.Shutdown(ctx)
  os.Exit(0)
}()
```
