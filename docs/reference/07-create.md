# Create records

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

### Create a record

```go
created, err := client.Post.CreateOne(
    // required fields
    Post.Title.Set("what up"),
    Post.Published.Set(true),

    // optional fields
    Post.ID.Set("id"),
    Post.Content.Set("stuff"),
).Exec(ctx)
```

### Create a record with a relation

Use the `Link` method to connect new records with existing ones. For example, the following query creates a new post and sets the postID attribute of the comment.

```go
created, err := client.Comment.CreateOne(
    Comment.Title.Set(title),
    Comment.Post.Link(
        Post.ID.Equals(postID),
    ),
    Comment.ID.Set("post"),
).Exec(ctx)
```

## Next steps

Learn how to [update data](08-update.md).
