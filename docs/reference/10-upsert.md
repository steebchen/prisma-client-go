# Upsert records

Use upsert to update or create records depending on whether it already exists or not.

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
```

### Upsert a record

Use UpsertOne to query for a document, define what to write when creating the document, and what to update if the document already exists.

```go
post, err := client.Post.UpsertOne(
    // query
    Post.ID.Equals("upsert"),
).Create(
    // set these fields if document doesn't exist already
    Post.Title.Set("title"),
    Post.Views.Set(0),
    Post.ID.Set("upsert"),
).Update(
    // update these fields if document already exists
    Post.Title.Set("new-title"),
    Post.Views.Increment(1),
).Exec(ctx)
if err != nil {
    panic(err)
}
```

## Next steps

Check out the details of [querying for relations](11-relations.md).
