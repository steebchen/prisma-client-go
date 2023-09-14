# Create records

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

  post   Post   @relation(fields: [postID], references: [id])
  postID String
}
```

### Create a record

```go
created, err := client.Post.CreateOne(
  // required fields
  db.Post.Published.Set(true),
  db.Post.Title.Set("what up"),

  // optional fields
  db.Post.ID.Set("id"),
  db.Post.Content.Set("stuff"),
).Exec(ctx)
```

### Create a record with a relation

Use the `Link` method to connect new records with existing ones. For example, the following query creates a new comment
and sets the postID attribute of the comment.

```go
created, err := client.Comment.CreateOne(
  db.Comment.Content.Set("content"),
  db.Comment.Post.Link(
    db.Post.ID.Equals("id"),
  ),
  db.Comment.ID.Set("post"),
).Exec(ctx)
```
