# Raw API

You can use the raw API when there's something you can't do with the current go client features. The query will be
redirected to the underlying database, so everything supported by the database should work. Please note that you need to
use the syntax specific to the database you're using.

The examples use the following prisma schema:

```prisma
model Post {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    published Boolean
    title     String
    content   String?

    comments Comment[]
}

model Comment {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    content   String

    post   Post @relation(fields: [postID], references: [id])
    postID String
}
```

## MySQL & SQLite

### Query

Use `QueryRaw` to query for data and automatically unmarshal it into a slice of structs.

#### Select all

```go
var posts []db.PostModel
err := client.Prisma.QueryRaw(`SELECT * FROM Post`).Exec(ctx, &posts)
```

#### Select specific

```go
var posts []PostModel
err := client.Prisma.QueryRaw(`SELECT * FROM Post WHERE id = ? AND title = ?`, "123abc", "my post").Exec(ctx, &posts)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `count`, which contains the affected rows.

```go
count, err := client.Prisma.ExecuteRaw(`UPDATE Post SET title = ? WHERE id = ?`, "my post", "123").Exec(ctx)
```

## Postgres

### Query

#### Select all

```go
var posts []PostModel
err := client.Prisma.QueryRaw(`SELECT * FROM "Post"`).Exec(ctx, &posts)
```

#### Select specific

```go
var posts []PostModel
err := client.Prisma.QueryRaw(`SELECT * FROM "Post" WHERE id = $1 AND title = $2`, "id2", "title2").Exec(ctx, &posts)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `count`, which contains the affected rows.

```go
count, err := client.Prisma.ExecuteRaw(`UPDATE "Post" SET title = $1 WHERE id = $2`, "my post", "123").Exec(ctx)
```

## Next steps

Ensure consistency with [transactions](13-transactions.md).
