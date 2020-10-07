# Create records

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
