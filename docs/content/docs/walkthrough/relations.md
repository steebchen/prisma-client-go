# Relations

The examples use the following prisma schema:

```prisma
model User {
  id    String @id @default(cuid())
  name  String
  posts Post[]
}

model Post {
  id        String   @id @default(cuid())
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  published Boolean
  title     String
  content   String?

  // optional author
  user   User?   @relation(fields: [userID], references: [id])
  userID String?

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

### Find by nested relation

In a query, you can query for relations by using "Some", "Every" or "None".

Please [see the caveats](https://github.com/prisma/prisma/issues/18193) for the "Every" filter.

```go
// get posts which have at least one comment with a content "My Content" and that post's titles are all "What up?"
posts, err := client.Post.FindMany(
  db.Post.Title.Equals("What up?"),
  db.Post.Comments.Some(
    db.Comment.Content.Equals("My Content"),
  ),
).Exec(ctx)
```

You can nest relation queries as deep as you like:

```go
users, err := client.User.FindMany(
  db.User.Name.Equals("Author"),
  db.User.Posts.Some(
    db.Post.Title.Equals("What up?"),
    db.Post.Comments.Some(
      db.Comment.Content.Equals("My Content"),
    ),
    db.Post.Comments.None(
      db.Comment.Content.Equals("missing"),
    ),
  ),
).Exec(ctx)
```
