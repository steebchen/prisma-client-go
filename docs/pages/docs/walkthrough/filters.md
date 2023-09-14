# Query filters

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

## Type filters

You probably want to build detailed queries, such as if a database column contains a word,
or if a number is great or lower than something. On this page, you can explore what filter methods are available
by type. All of these queries are fully type-safe and independent of the underlying database.

### String filters

```go
// query for posts where the title is "my post"
db.Post.Title.Equals("my post"),
// query for titles containing the string "post"
db.Post.Title.Contains("post"),
// query for titles starting with "my"
db.Post.Title.StartsWith("my"),
// query for titles ending with "post"
db.Post.Title.EndsWith("post"),
```

### Number filters

```go
// query for all posts which have exactly 50 views
db.Post.Views.Equals(50),
// query for all posts which have less than or exactly 50 views
db.Post.Views.Lte(50),
// query for all posts which have less than 50 views
db.Post.Views.Lt(50),
// query for all posts which have more than or exactly 50 views
db.Post.Views.Gte(50),
// query for all posts which have more than 50 views
db.Post.Views.Gt(50),
```

### Time filters

```go
// query for all posts which equal an exact date
db.Post.CreatedAt.Equals(yesterday),
// query for all posts which were created in the last 6 hours(createdAt > 6 hours ago)
db.Post.CreatedAt.Gt(time.Now().Add(-6 * time.Hour)),
// query for all posts which were created in the last 6 hours(createdAt >= 6 hours ago)
db.Post.CreatedAt.Gte(time.Now().Add(-6 * time.Hour)),
// query for all posts which were created until yesterday
db.Post.CreatedAt.Lt(time.Now().Truncate(24 * time.Hour)),
// query for all posts which were created until yesterday including today's 00:00:00
db.Post.CreatedAt.Lte(time.Now().Truncate(24 * time.Hour)),
```

### Optional type filters

Optional fields are hard to represent in Go, since SQL has NULLs but Go does not have nullable types.
Usually, the community defaults to using pointers, but providing that everywhere can be inconvenient. In order to set
NULLs by using a pointer, you can use the `XOptional` method variants.

```go
// set an optional field with a specific string
db.Post.Content.Equals("my description")

// set an optional field by using a pointer, where a nil pointer means
// to set NULL in the database
db.Post.Content.EqualsOptional(nil)

// or by using a pointer
content := "string"
// ...
db.Post.Content.EqualsOptional(&content)
```

## General

There are a few general filters you can apply. Note that the model has to be used to preserve type information.

### Not

If you want to negate a query, you can use `Not`.

The following query queries for all posts where their title doesn't equal "123":

```go
db.Post.Not(
  db.Post.Title.Equals("123"),
)
```

### Or

You can use `Or`.

The following query queries for all posts where either their title equals "123" OR their content equals "456":

```go
db.Post.Or(
  db.Post.Title.Equals("123"),
  db.Post.Content.Equals("456"),
)
```
