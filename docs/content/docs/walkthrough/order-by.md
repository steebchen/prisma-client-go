# Order By

The examples use the following prisma schema:

```prisma
model Post {
  id        String   @id @default(cuid())
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  published Boolean
  title     String
  content   String?

  // add an index to be able to order by created_at
  @@index([createdAt])
}
```

### Order by ID

The following example would equal to the default behaviour of ordering by ID in ascending order:

```go
posts, err := client.Post.FindMany().OrderBy(
  db.Post.ID.Order(db.SortOrderAsc),
).Exec(ctx)
```

You can order by any field ein either direction, but it's recommended to use an index on fields you order.

#### Order by latest created

```go
posts, err := client.Post.FindMany().OrderBy(
  db.Post.CreatedAt.Order(db.SortOrderDesc),
).Exec(ctx)
```

#### Combine with pagination

```go
posts, err := client.
  Post.
  FindMany().
  Take(5).
  Cursor(db.Post.ID.Cursor("abc")).
  OrderBy(
    db.Post.CreatedAt.Order(db.SortOrderDesc),
  ).Exec(ctx)
```
