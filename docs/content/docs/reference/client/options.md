# Client options

Configuration or options in the Prisma client are set via functional options.

## WithDatasourceURL

You can configure or override the datasource URL at runtime:

```go
client := db.NewClient(
  db.WithDatasourceURL("postgresql://localhost:5432/mydb?schema=public"),
)
```
