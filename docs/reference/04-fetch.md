# Fetching additional data

You can query for an entity and specify what to return in addition. For example, if you want to show a post's information with some of its comments, you would usually do 2 separate queries, but using the With/Fetch syntax you can do it in a single query.

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

### Find a post and fetch three of its comments

```go
// find a post
post, err := client.Post.FindUnique(
    Post.Title.Equals("hi"),
).With(
    // also fetch 3 of its comments
    Post.Comments.Fetch().Take(3),
).Exec(ctx)
check(err)
log.Printf("post's title: %s", post.Title)

comments := post.Comments()
for _, comment := range comments {
    log.Printf("comment: %s", comment)
}
```

## Next steps

Check out a [detailed explanation for creating rows](07-create.md).
