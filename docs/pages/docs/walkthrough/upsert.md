# Upsert records

Use upsert to update or create records depending on whether it already exists or not.

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
}
```

### Upsert a record

Use UpsertOne to query for a document, define what to write when creating the document, and what to update if the
document already exists.

```go
post, err := client.Post.UpsertOne(
  // query
  db.Post.ID.Equals("upsert"),
).Create(
  // set these fields if document doesn't exist already
  db.Post.Published.Set(true),
  db.Post.Title.Set("title"),
  db.Post.ID.Set("upsert"),
).Update(
  // update these fields if document already exists
  db.Post.Title.Set("new-title"),
  db.Post.Views.Increment(1),
).Exec(ctx)
```
