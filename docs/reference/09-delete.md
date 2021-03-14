# Delete records

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

### Delete a record

To delete a record, just query for a field using FindUnique or FindMany, and then just chain it by invoking `.Delete()`.

```go
deleted, err := client.Post.FindUnique(
    Post.Title.Equals("what up"),
).Delete().Exec(ctx)
```

## Next steps

Learn how to [upsert documents](10-upsert.md).
