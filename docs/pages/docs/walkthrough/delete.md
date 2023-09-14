# Delete records

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

### Delete a record

To delete a record, just query for a field using FindUnique or FindMany, and then just chain it by invoking `.Delete()`.

```go
deleted, err := client.Post.FindUnique(
  db.Post.ID.Equals("id"),
).Delete().Exec(ctx)
```
