# Basic API

Find, update and delete records.

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

## Reading data

### Find many records

```go
posts, err := client.Post.FindMany(
    db.Post.Title.Equals("hi"),
).Exec(ctx)
```

If no records are found, this returns an empty array without returning an error (like usual SQL queries).

### Find one record

```go
post, err := client.Post.FindOne(
    db.Post.ID.Equals("123"),
).Exec(ctx)

if err == db.ErrNotFound {
    log.Printf("no record with id 123")
}
```

This returns an error of type `ErrNotFound` (exported in the `db` package) if there was no such record.

### Query API

Depending on the data types of your fields, you will automatically be able to query for respective operations. For example, for integer or float fields you might want to query for a field which is less than or greater than some number.

```go
post, err := client.Post.FindOne(
    // query for posts containing the title "hi"
    db.Post.Title.Contains("what up"),
).Exec(ctx)
```

To explore more query filter, see [all possible query filters](./03-filters.md).

### Querying for relations

In a query, you can query for relations by using "Some" or "Every". You can also query for deeply nested relations.

```go
// get posts which have at least one comment with a title "My Title" and that post's comments are all "What up?"
posts, err := client.Post.FindMany(
    Post.Title.Equals("what up"),
    Post.Comments.Some(
        Comment.Title.Equals("My Title"),
    ),
).Exec(ctx)
```

To explore querying for relations in detail, see [more relation query examples](./08-relations.md).

## Next steps

Read the next article [query filters](./03-filters.md) to explore how to form more complex queries.
