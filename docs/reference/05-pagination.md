# Pagination

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

### Return the first 5 rows

```go
created, err := client.
    Post.
    FindMany().
    Take(5).
    Exec(ctx)
```

### Return the first 5 rows and skip 2 rows

```go
created, err := client.
    Post.
    FindMany().
    Take(5).
    Skip(2).
    Exec(ctx)
```

## Cursor-based pagination

Instead of using `Skip`, you can also provide a cursor:

```go
created, err := client.
    Post.
    FindMany().
    Take(5).
    Skip(2).
    Cursor(Post.ID.Cursor("abc")).
    Exec(ctx)
```

Also check out the [order by docs](06-order-by.md) to understand how you can combine cursor-based pagination with order by.

## Next steps

Learn how to [order queries](07-create.md).
