# Transactions

A database transaction refers to a sequence of read/write operations that are guaranteed to either succeed or fail as a whole.

## Successful scenario

A simple transaction could look as follows. Just omit the `Exec(ctx)`, and provide the Prisma calls to `client.Prisma.Transaction`:

```go
// create two users at once and run in a transaction

createUserA := client.User.CreateOne(
    User.Email.Set("a"),
    User.ID.Set("a"),
)

createUserB := client.User.CreateOne(
    User.Email.Set("b"),
    User.ID.Set("b"),
)

if err := client.Prisma.Transaction(createUserA, createUserB).Exec(ctx); err != nil {
    panic(err)
}
```

## Failure scenario

Let's say we have one user record in the database:

```json
{
    "id": "123",
    "email": "john@example.com"
}
```

```go
// this will fail, since the record doesn't exist...
a := client.User.FindUnique(
    User.ID.Equals("does-not-exist"),
).Update(
    User.Email.Set("foo"),
)

// ...so this should be roll-backed, even though itself it would succeed
b := client.User.FindUnique(
    User.ID.Equals("123"),
).Update(
    User.Email.Set("new-email@doe.com"),
)

if err := client.Prisma.Transaction(a, b).Exec(ctx); err != nil {
    // this err will be non-nil and the transaction will rollback,
    // so nothing will be updated in the database
    panic(err)
}
```

## Next steps

Check out how to use [json fields](13-json.md).
