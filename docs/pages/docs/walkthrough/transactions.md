# Transactions

A database transaction refers to a sequence of read/write operations that are guaranteed to either succeed or fail as a
whole.

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

## Example

A simple transaction could look as follows. Just omit the `Exec(ctx)`, and provide the Prisma calls
to `client.Prisma.Transaction`:

```go
// create two posts at once and run in a transaction

firstPost := client.Post.CreateOne(
  db.Post.Published.Set(true),
  db.Post.Title.Set("First Post"),
).Tx()

secondPost := client.Post.CreateOne(
  db.Post.Published.Set(false),
  db.Post.Title.Set("Second Post"),
).Tx()

if err := client.Prisma.Transaction(firstPost, secondPost).Exec(ctx); err != nil {
  panic(err)
}

log.Printf("first post result: %+v", firstPost.Result())
log.Printf("second post result: %+v", secondPost.Result())
```

## Multiple transactions

If you have many transaction items or need to construct a transaction dynamically, you can use the `PrismaTransaction` type:

```go
// create two posts at once and run in a transaction
var txns []db.PrismaTransaction

for _, title := range []string{"First Post", "Second Post"} {
  txn := client.Post.CreateOne(
    db.Post.Published.Set(true),
    db.Post.Title.Set(title),
  ).Tx()
  txns = append(txns, txn)
}

if err := client.Prisma.Transaction(txns...).Exec(ctx); err != nil {
  panic(err)
}
```

## Setting relations with foreign keys

When setting relations with foreign keys, you have to set the IDs of the related records. For example, to create a post with multiple comments, you would do the following:

```go
postID := "12345"

txns := []db.PrismaTransaction{
  client.Post.CreateOne(
    db.Post.Published.Set(true),
    db.Post.Title.Set("First Post"),
    db.Post.ID.Set(postID),
  ).Tx(),
  client.Comment.CreateOne(
    db.Comment.Content.Set("First comment"),
    db.Comment.Post.Link(
      db.Post.ID.Equals(postID),
    ),
  ).Tx(),
  client.Comment.CreateOne(
    db.Comment.Content.Set("Second comment"),
    db.Comment.Post.Link(
      db.Post.ID.Equals(postID),
    ),
  ).Tx(),
}

if err := client.Prisma.Transaction(txns...).Exec(ctx); err != nil {
  panic(err)
}
```

## Failure scenario

Let's say we have one post record in the database:

```json
{
  "id": "123",
  "title": "Hi from Prisma"
}
```

```go
// this will fail, since the record doesn't exist...
a := client.Post.FindUnique(
  db.Post.ID.Equals("does-not-exist"),
).Update(
  db.Post.Title.Set("new title"),
).Tx()

// ...so this should be roll-backed, even though itself it would succeed
b := client.Post.FindUnique(
  db.Post.ID.Equals("123"),
).Update(
  db.Post.Title.Set("New title"),
).Tx()

if err := client.Prisma.Transaction(b, a).Exec(ctx); err != nil {
  // this err will be non-nil and the transaction will rollback,
  // so nothing will be updated in the database
  panic(err)
}
```
