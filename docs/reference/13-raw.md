# Raw API

You can use the raw API when there's something you can't do with the current go client features. The query will be
redirected to the underlying database, so everything supported by the database should work. Please note that you need to
use the syntax specific to the database you're using.

NOTE: When defining your return type structure, you have to use Prisma-specific raw data types,
such as `RawInt`, `RawString`, etc. exported by the Prisma client due to how the Prisma internals
and Go works. If you are querying for a specific model, you can also use `Raw<Model>`, e.g. `RawUser`
instead of `UserModel`.

The examples use the following prisma schema:

```prisma
model Post {
    id        String   @id @default(cuid())
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    published Boolean
    title     String
    content   String?
    views     Int      @default(0)

    comments Comment[]
}

model Comment {
    id        String   @id @default(cuid())
    createdAt DateTime @default(now())
    content   String

    post   Post   @relation(fields: [postID], references: [id])
    postID String
}
```

## MySQL & SQLite

### Query

Use `QueryRaw` to query for data and automatically unmarshal it into a slice of structs.

#### Select all for a model

```go
var posts []db.RawPost
err := client.Prisma.QueryRaw(`SELECT * FROM `Post``).Exec(ctx, &posts)
```

#### Select specific

```go
// note the usage of RawPost instead of PostModel
var posts []db.RawPost
err := client.Prisma.QueryRaw("SELECT * FROM `Post` WHERE id = ? AND title = ?", "123abc", "my post").Exec(ctx, &posts)
```

### Custom Query

```go
// note the usage of db.RawString, db.RawInt, etc.
var res []struct{
	PostID   db.RawString `json:"post_id"`
	Comments db.RawInt    `json:"comments"`
}
err := client.Prisma.QueryRaw("SELECT post_id, count(*) as comments FROM `Comment` GROUP BY post_id").Exec(ctx, &res)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `result.Count`, which contains the affected rows.

```go
result, err := client.Prisma.ExecuteRaw("UPDATE `Post` SET title = ? WHERE id = ?", "my post", "123").Exec(ctx)
println(result.Count) // 1
```

## Postgres

### Query

#### Select all for a model

```go
var posts []db.RawPost
err := client.Prisma.QueryRaw(`SELECT * FROM "Post"`).Exec(ctx, &posts)
```

#### Select specific

```go
var posts []db.RawPost
err := client.Prisma.QueryRaw(`SELECT * FROM "Post" WHERE id = $1 AND title = $2`, "id2", "title2").Exec(ctx, &posts)
```

### Custom Query

```go
// note the usage of db.RawString, db.RawInt, etc.
var res []struct{
	PostID   db.RawString `json:"post_id"`
	Comments db.RawInt    `json:"comments"`
}
err := client.Prisma.QueryRaw(`SELECT post_id, count(*) as comments FROM "Comment" GROUP BY post_id`).Exec(ctx, &res)
```

### Operations

Use `ExecuteRaw` for operations such as `INSERT`, `UPDATE` or `DELETE`. It will always return a `count`, which contains the affected rows.

```go
result, err := client.Prisma.ExecuteRaw(`UPDATE "Post" SET title = $1 WHERE id = $2`, "my post", "123").Exec(ctx)
println(result.Count) // 1
```

## Next steps

Ensure consistency with [transactions](14-transactions.md).
