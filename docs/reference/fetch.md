# Fetching additional data

You can query for an entity and specify what to return in addition. For example, if you want to show a user's information with some of their posts, you would usually do 2 separate queries, but using the With/Fetch syntax you can do it in a single query.

```go
// find a user
user, err := client.User.FindOne(
    User.Email.Equals("john@example.com"),
).With(
    // also fetch 3 their posts
    User.Posts.Fetch().Take(3),
).Exec(ctx)
check(err)
log.Printf("user's name: %s", user.Name)

posts := user.Posts()
for _, post := range posts {
    log.Printf("post: %+v", post)
}
```
