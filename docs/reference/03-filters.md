# Query filters

## Type filters

You probably want to build detailed queries, such as if a database column contains a word,
or if a number is great or lower than something. On this page, you can explore what filter methods are available
by type. All of these queries are fully type-safe and independent of the underlying database.

### String filters

```go
// query for people who are named "John"
db.Post.Title.Equals("John"),
// query for names containing the string "oh"
db.Post.Title.Contains("oh"),
// query for names starting with "Jo"
db.Post.Title.HasPrefix("Jo"),
// query for names ending with "Jo"
db.Post.Title.HasSuffix("hn"),
```

### Number filters

```go
// query for all posts which have exactly 50 views
db.Post.Views.Equals(50),
// query for all posts which have less than or exactly 50 views
db.Post.Views.LTE(50),
// query for all posts which have less than 50 views
db.Post.Views.LT(50),
// query for all posts which have more than or exactly 50 views
db.Post.Views.GT(50),
// query for all posts which have more than 50 views
db.Post.Views.GTE(50),
```


### Time filters

```go
// query for all posts which equal an exact date
db.Post.CreatedAt.Equals(yesterday),
// query for all posts which were created in the last 6 hours
db.Post.CreatedAt.After(time.Now().Add(-6 * time.Hour)),
// query for all posts which were created in the last 6 hours including right now
db.Post.CreatedAt.BeforeEquals(time.Now().Add(-6 * time.Hour)),
// query for all posts which were created until yesterday
db.Post.CreatedAt.Before(time.Now().Truncate(24 * time.Hour)),
// query for all posts which were created until yesterday including right now
db.Post.CreatedAt.BeforeEquals(time.Now().Truncate(24 * time.Hour)),
```

## General

There are a few general filters you can apply. Note that the model has to be used to preserve type information.

### Not

If you want to negate a query, you can use `Not`.

The following query queries for all posts where their title doesn't equal "123":

```go
db.Post.Not(
  db.Post.Title.Equals("123"),
)
```

### Or

If you want to negate a query, you can use `Or`.

The following query queries for all posts where either their title equals "123" OR their description equals "456":

```go
db.Post.Or(
  db.Post.Title.Equals("123"),
  db.Post.Desc.Equals("456"),
)
```

## Next steps

In the next article, you can explore how to [fetch for multiple things](./04-fetch.md) at once.
