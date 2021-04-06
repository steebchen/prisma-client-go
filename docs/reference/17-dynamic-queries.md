# Dynamic Queries

For more complex queries, you can build up type-safe queries by using exported interfaces.

The examples use the following prisma schema:

```prisma
model User {
    id       String   @default(cuid()) @id
    kind     String   // can be of type employee or customer
    email    String
    referrer String?
}

```

## Param

The params are exported the following shape:

```
<Model><Action>Param
```

For the model user above, there are two main exported interfaces, one for querying for data and one for writing data respectively:

```
UserWhereParam
UserSetParam
```

## Example

With the schema above and the CreateUser function users of type customers should have their IP address saved, but not for new users of type employee.

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var params []db.UserSetParam
    email := r.PostFormValue("email")
    kind := r.PostFormValue("kind")
    if kind == "customer" {
        // Set the referer for users of type customer only
        params = append(params, db.User.Referrer.Set("Referer"))
    }
    _, err := client.User.CreateOne(
        db.User.Kind.Set(kind),
        db.User.Email.Set(email),
        params...,
    ).Exec(r.Context())
    if err != nil {
        panic(err)
    }
    // write results to response
    // ...
}
```

## Next steps

The go client is still in an early access and has [limitations](18-limitations.md).
