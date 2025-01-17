# Field selection and omission

You can select or omit fields in the query API. This is useful for reducing the amount of data you need to fetch, or
for reducing the amount of data you need to send from the database to the client (up until to the end user)

The examples use the following prisma schema:

```prisma
model User {
  id        String   @id @default(cuid())
  name      String?
  password  String?
  age       Int?
}
```

## Notes

You can only select or omit fields in a query, not both.

## Select

Select returns only the fields you specify, and nothing else.

```go
users, err := client.User.FindMany(
  User.Name.Equals("john"),
).Select(
  User.Name.Field(),
).Exec(ctx)
if err != nil {
  panic(err)
}
```

## Omit

Omit returns all fields except the ones you specify.

```go
users, err := client.User.FindMany(
  User.Name.Equals("a"),
).Omit(
  User.ID.Field(),
  User.Password.Field(),
  User.Age.Field(),
  User.Name.Field(),
).Exec(ctx)
if err != nil {
  panic(err)
}
```
