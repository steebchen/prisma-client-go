# API reference

Contents

- [Introduction](01-introduction.md)

- [Find rows](02-find.md)
- [Query filters](03-filters.md)
- [Fetching multiple things at once](04-fetch.md)
- [Limit & Pagination](05-pagination.md)
- [Order by](06-order-by.md)

- [Create rows](07-create.md)
- [Update rows](08-update.md)
- [Delete rows](09-delete.md)

- [Query for relations](10-relations.md)

- [Raw API fallback](11-raw.md)
- [JSON](12-json.md)
- [Limitations](15-limitations.md)

```go
func main() {
    users, err := client.User.FindUnique().Exec()
}
```
