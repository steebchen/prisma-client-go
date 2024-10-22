# Update records

The examples use the following prisma schema:

```prisma
model Post {
  id        String   @id @default(cuid())
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  published Boolean
  title     String
  content   String?

  comments Comment[]
}

model Comment {
  id        String   @id @default(cuid())
  createdAt DateTime @default(now())
  content   String

  post   Post?   @relation(fields: [postID], references: [id])
  postID String?
}
```

### Update a record

To update a record, just query for a field using FindUnique or FindMany, and then just chain it by invoking `.Update()`.

```go
updated, err := client.Post.FindMany(
  db.Post.Title.Equals("what up"),
).Update(
  db.Post.Content.Set("new content"),
  db.Post.Title.Set("new title"),
).Exec(ctx)
```

### Update relations

#### Required relation

You can set relations in the same way as when creating records.

```go
updated, err := client.Comment.FindUnique(
  db.Comment.ID.Equals("id"),
).Update(
  db.Comment.Post.Link(
    db.Post.ID.Equals(postID),
  ),
).Exec(ctx)
```

#### Optional relation

For optional relations, you can also unlink the relation, so the foreign key value is set to `NULL`:

```go
updated, err := client.Comment.FindUnique(
  db.Comment.ID.Equals("id"),
).Update(
  db.Comment.Post.Unlink(),
).Exec(ctx)
```
