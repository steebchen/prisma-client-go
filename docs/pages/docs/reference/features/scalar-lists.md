# Scalar lists

The examples use the following prisma schema:

```prisma
model Post {
  id    String @id @default(cuid())
  title String
  tags  String[]
}
```

### Querying records by array fields

Note: This only works with postgres.

```go
post, err := client.Post.FindFirst(
  // whether the list contains a single field
  db.Post.Tags.Has("coffee"),
  // or
  // whether the tag contains coffee _and_ juice
  db.Post.Tags.HasEvery([]string{"coffee", "juice"}),
  // or
  // whether the tag contains coffee or tea
  db.Post.Tags.HasSome([]string{"coffee", "tea"}),
  // or
  db.Post.Tags.IsEmpty(false),
).Exec(ctx)
```

### Writing to array fields

Set the field regardless of the previous value:

```go
post, err := client.Post.FindUnique(
  db.Post.ID.Equals("123"),
).Update(
  db.Post.Tags.Set([]string{"a", "b", "c"}),
).Exec(ctx)
```

Add items to an existing list (postgres only):

```go
post, err := client.Post.FindUnique(
  db.Post.ID.Equals("123"),
).Update(
  db.Post.Tags.Push([]string{"a", "b"}),
).Exec(ctx)
```

### Notes

NULL values in scalar
lists [may need extra consideration](https://www.prisma.io/docs/concepts/components/prisma-client/working-with-fields/working-with-scalar-lists-arrays#filtering-scalar-lists).
