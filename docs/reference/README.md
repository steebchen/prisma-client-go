# API reference

Contents

- [Introduction](./01-introduction.md)

- [Find rows](./02-find.md)
- [Query filters](./03-filters.md)
- [Fetching multiple things at once](./04-fetch.md)

- [Create rows](./05-create.md)
- [Update rows](./06-update.md)
- [Delete rows](./07-delete.md)

- [Query for relations](./08-relations.md)

- [Raw API fallback](./09-raw.md)
- [Limitations](./10-limitations.md)

```go
func main() {
  users, err := client.User.FindOne().Exec()
}
```
