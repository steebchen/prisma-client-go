# Order By

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

### Order by ID

The following example would equal to the default behaviour of ordering by ID in ascending order:

```go
posts, err := client.Post.FindMany().OrderBy(
    Post.ID.Order(ASC),
).Exec(ctx)
```

You can order by any field ein either direction, but it's recommended to use an index on fields you order.

#### Order by latest created

```go
posts, err := client.Post.FindMany().OrderBy(
    Post.CreatedAt.Order(DESC),
).Exec(ctx)
```

#### Combine with pagination

```go
posts, err := client.Post.FindMany().Take(5).Cursor(Post.CreatedAt.Cursor(someDate)).OrderBy(
    Post.CreatedAt.Order(DESC),
).Exec(ctx)
```

## Next steps

Check out a [detailed explanation on how to create records](07-create.md).
