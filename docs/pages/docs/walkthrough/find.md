# Basic API

Find, update and delete records.

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

## Reading data

### Find many records

```go
posts, err := client.Post.FindMany(
  db.Post.Title.Equals("hi"),
).Exec(ctx)
```

If no records are found, the query above returns a slice without returning an error (like normal SQL queries).

### Find a unique record

FindUnique finds a record which is guaranteed to be unique, like @id fields or fields marked with @unique.

```go
post, err := client.Post.FindUnique(
  db.Post.ID.Equals("123"),
).Exec(ctx)
if errors.Is(err, db.ErrNotFound) {
  log.Printf("no record with id 123")
} else if err != nil {
  log.Printf("error occurred: %s", err)
}
```

### Find a single record

FindFirst finds the first record found. It has the same query capabilities as FindMany, but acts as a convenience method
to return just the first record found.

```go
post, err := client.Post.FindFirst(
  db.Post.Title.Equals("hi"),
).Exec(ctx)
if errors.Is(err, db.ErrNotFound) {
  log.Printf("no record with title 'hi' found")
} else if err != nil {
  log.Printf("error occurred: %s", err)
}

log.Printf("post: %+v", post)
```

This returns an `ErrNotFound` error (exported by the generated client) if there was no such record.

### Query API

The query operations change based on the data types in your schema. For example, integers and floats will have greater
than and less than operations, while strings have prefix and suffix operations.

```go
posts, err := client.Post.FindMany(
  // query for posts containing the title "What"
  db.Post.Title.Contains("What"),
).Exec(ctx)
```

To explore more query filters, see [all possible query filters](filters.md).

### Querying for relations

You can query for relations by using "Some", "Every" or "None" to query for records where only some, all or none of the records match
respectively. You can nest those queries as deep as you like.

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

To explore querying for relations in detail, see [more relation query examples](relations.md).
