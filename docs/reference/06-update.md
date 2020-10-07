# Update records

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

### Update a record

To update a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Update()`.

```go
updated, err := client.Post.FindOne(
    Post.Title.Equals("what up"),
).Update(
    Post.Desc.Set("new description"),
    Post.Title.Set("new title"),
).Exec(ctx)
```

## Next steps

Learn how to [delete data](./07-delete.md).
