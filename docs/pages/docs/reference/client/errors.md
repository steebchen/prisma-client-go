# Errors

## ErrNotFound

`ErrNotFound` is returned when a query does not return any results. This error may be returned in `FindUnique`, `FindFirst`, but also when updating or deleting single records using `FindUnique().Update()` and `FindUnique().Delete()`.

```go
post, err := client.Post.FindFirst(
  db.Post.Title.Equals("hi"),
).Exec(ctx)
if err != nil {
  if errors.Is(err, db.ErrNotFound) {
    panic("no record with title 'hi' found")
  }
  panic("error occurred: %s", err)
}
```

## IsErrUniqueConstraint

A unique constraint violation happens when a query attempts to insert or update a record with a value that already exists in the database, or in other words, violates a unique constraint.

```go
user, err := db.User.CreateOne(...).Exec(cxt)
if err != nil {
  if info, err := db.IsErrUniqueConstraint(err); err != nil {
    // Fields exists for Postgres and SQLite
    log.Printf("unique constraint on the fields: %s", info.Fields)

    // you can also compare it with generated field names:
    if info.Fields[0] == db.User.Name.Field() {
      // do something
      log.Printf("unique constraint on the `user.name` field")
    }

    // For MySQL and MongoDB, use the constraint key
    log.Printf("unique constraint on the key: %s", info.Key)
  }
}
```
