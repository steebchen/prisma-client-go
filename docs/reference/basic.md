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

Other possible queries are:

```go
// query for people who are named "John"
db.Post.Title.Contains("John"),
// query for names containing the string "oh"
db.Post.Title.Contains("oh"),
// query for names starting with "Jo"
db.Post.Title.HasPrefix("Jo"),
// query for names ending with "Jo"
db.Post.Title.HasSuffix("hn"),
// query for all posts which have less than or exactly 50 views
db.Post.Views.LTE(50),
// query for all posts which have less than 50 views
db.Post.Views.LT(50),
// query for all posts which have more than or exactly 50 views
db.Post.Views.GT(50),
// query for all posts which have more than 50 views
db.Post.Views.GTE(50),
// query for all posts which were created in the last 6 hours
db.Post.CreatedAt.After(time.Now().Add(-6 * time.Hour)),
// query for all posts which were created until yesterday
db.Post.CreatedAt.Before(time.Now().Truncate(24 * time.Hour)),
```

All of these queries are fully type-safe and independent of the underlying database.

### Querying for relations

In a query, you can query for relations by using "Some" or "Every". You can also query for deeply nested relations.

```go
// get a post which has at least one comment with a title "My Title" and that post's comments are all "What up?"
actual, err := client.Post.FindMany(
    Post.Title.Equals("what up"),
    Post.Comments.Some(
        Comment.Title.Equals("My Title"),
    ),
).Exec(ctx)
```

## Writing data

### Create a record

```go
created, err := client.Post.CreateOne(
    // required fields
    Post.Title.Set("what up"),
    Post.Desc.Set("this is a description"),

    // optional fields
    Post.ID.Set("id"),
    Post.Title.Set("name"),
    Post.Stuff.Set("stuff"),
).Exec(ctx)
```

### Create a record with a relation

Use the method `Link` to connect new objects with existing ones. For example, the following query creates a new post and sets the postID attribute of the comment.

```go
created, err := client.Comment.CreateOne(
    Comment.Title.Set(title),
    Comment.Post.Link(
        Post.ID.Equals(postID),
    ),
    Comment.ID.Set("post"),
).Exec(ctx)
```

### Update a record

To update a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Update()`.

```go
updated, err := client.Post.FindOne(
    Post.Title.Equals("what up"),
).Update(
    Post.Desc.Set("new description"),
    Post.Title.Set("new title"),
).Exec(ctx)
```

### Delete a record

To delete a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Delete()`.

```go
updated, err := client.Post.FindOne(
    Post.Title.Equals("what up"),
).Delete().Exec(ctx)
```
