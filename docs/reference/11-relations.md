# Relations

The examples use the following prisma schema:

```prisma
model User {
    id    String   @default(cuid()) @id
    name  String
    posts Post[]
}

model Post {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    published Boolean
    title     String
    content   String?

    // optional author
    user   User @relation(fields: [userID], references: [id])
    userID String

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

### Find by nested relation

In a query, you can query for relations by using "Some" or "Every":

```go
// get posts which have at least one comment with a title "My Title" and that post's comments are all "What up?"
posts, err := client.Post.FindMany(
    Post.Title.Equals("what up"),
    Post.Comments.Some(
        Comment.Title.Equals("My Title"),
    ),
).Exec(ctx)
```

You can nest relation queries as deep as you like:

```go
users, err := client.User.FindMany(
    Post.Title.Equals("what up"),
    Post.Comments.Some(
        Comment.Title.Equals("My Title"),
    ),
).Exec(ctx)
```

## Next steps

If the Go client shouldn't support for a query you need to do, read how you can use [raw SQL queries](12-raw.md).
