# Raw API

You can use the raw API when there's something you can't do with the current go client features. The query will be
redirected to the underlying database, so everything supported by the database should work. Please note that you need to
use the syntax specific to the database you're using.

## MySQL & SQLite

### Query

Use `QueryRaw` to query for data and automatically unmarshal it into a user-defined slice.

#### Select all

```go
var users []db.UserModel
err := client.QueryRaw(`SELECT * FROM User`).Exec(ctx, &users)
```

#### Select specific

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM User WHERE id = ? AND email = ?`, "123abc", "prisma@example.com").Exec(ctx, &users)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `count`, which contains the affected rows.

```go
count, err := client.ExecuteRaw(`UPDATE User SET name = ? WHERE id = ?`, "John", "123").Exec(ctx)
```

## Postgres

### Query

#### Select all

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM "User"`).Exec(ctx, &users)
```

#### Select specific

```go
var users []UserModel
err := client.QueryRaw(`SELECT * FROM "User" WHERE id = $1 AND email = $2`, "id2", "email2").Exec(ctx, &users)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `count`, which contains the affected rows.

```go
count, err := client.ExecuteRaw(`UPDATE "User" SET name = $1 WHERE id = $2`, "John", "123").Exec(ctx)
```

## Next steps

Check out how to use [json fields](12-json.md).
