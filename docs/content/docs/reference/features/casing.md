# Case sensitivity

Note that case sensitivity depends on the database and database collation.

MySQL uses case insensitivity by default, while postgres is case-sensitive.

It's recommended to check out
the [full docs about case sensitivity](https://www.prisma.io/docs/concepts/components/prisma-client/case-sensitivity).

## Explicitly query for case (in)sensitive data

This is only generated for Postgres.

```go
users, err := client.User.FindMany(
User.Email.Equals("prisMa"),
User.Email.Mode(QueryModeInsensitive), // sets case insensitivity
).Exec(ctx)
```

```
+----+-----------------------------------+
| id | email                             |
+----+-----------------------------------+
| 61 | alice@prisma.io                   |
| 49 | brigitte@prisma.io                |
+----+-----------------------------------+
```
