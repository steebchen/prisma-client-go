# XIfPresent method variants

If you want to query for data and dynamically determine whether a field should be ignored or not, or if you want expose
an update operation where only some fields may get set, you can use the IfPresent method variants.

This does NOT mean SQL NULL – use XOptional method variants for that instead.

The examples use the following prisma schema:

```prisma
model Post {
  id      String  @id @default(cuid())
  title   String
  content String?
}
```

## Querying dynamically

You might want to query dynamically if you have an API and you want the end-user to decide which fields to query. In the
following example, the fields title and content are queried, but if a variable is nil, it means the field should be
ignored.

```go
title := "hi"
var content *string
_, err := client.Post.FindMany(
  // query for this one
  db.Post.Title.EqualsIfPresent(&title),
  // ignore this one, since `content` nil
  db.Post.Content.EqualsIfPresent(content),
).Exec(ctx)
```

## Writing data dynamically

Writing data dynamically works the same way as querying. If a pointer is nil, the field will not be touched; if it's
present, the field value will be updated.

```go
var newTitle *string
newContent := "hi"
_, err := client.Post.FindUnique(
  db.Post.ID.Equals("123"),
).Update(
  // don't set because `newTitle` is nil
  db.Post.Title.SetIfPresent(newTitle),
  // set value
  db.Post.Content.SetIfPresent(&newContent),
).Exec(ctx)
```
