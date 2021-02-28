# Mocks for testing

When writing tests for functions which invoke methods used by the Go client, you usually need to run a real database in order to make these tests work. While this is acceptable for integration tests, it may be harder to test code because you need to reset and migrate the database, seed data, and then tear off everything. Plus, it results in incredibly slow test runs.

Prisma provides native mocks in order to inject anything you want in your unit tests. This way, you can define for what query what result is returned, and then test your function with that information. They will also be fast, because there's no real database needed, and you just define yourself what a function should return.

The examples use the following prisma schema:

```prisma
model Post {
    id    String   @default(cuid()) @id
    title String
}
```

## Testing for results

The function `GetPostTitle` acts as your resolver; for example, it could be a function which is called in an API route.

To write a unit test, you would create a new Prisma mock client, define your expectations, and then run your actual test.
Expectations consist of the exact query or queries you expect, and the result what should be returned.

```go
// main.go
func GetPostTitle(ctx context.Context, client *PrismaClient, postID string) (string, error) {
    post, err := client.Post.FindUnique(
        db.Post.ID.Equals(postID),
    ).Exec(ctx)
    if err != nil {
        return "", fmt.Errorf("error fetching post: %w", err)
    }

    return post.Title, nil
}

// main_test.go
func TestGetPostTitle_returns(t *testing.T) {
    // create a new mock
    // this returns a mock prisma `client` and a `mock` object to set expectations
    client, mock, ensure := NewMock()
    // defer calling ensure, which makes sure all of the expectations were met and actually called
    // calling this makes sure that an error is returned if there was no query happening for a given expectation
    // and makes sure that all of them succeeded
    defer ensure(t)

    expected := &PostModel{
        InnerPost: InnerPost{
            ID:   "123",
            Title: "foo",
        },
    }

    // start the expectation
    mock.Post.Expect(
        // define your exact query as in your tested function
        // call it with the exact arguments which you expect the function to be called with
        // you can copy and paste this from your tested function, and just put specific values into the arguments
        client.Post.FindUnique(
            db.Post.ID.Equals("123"),
        ),
    ).Returns(expected) // sets the object which should be returned in the function call

    // mocking set up is done; let's define the actual test now
    title, err := GetPostTitle(context.Background(), client, "123")
    if err != nil {
        t.Fatal(err)
    }

    if title != "foo" {
        t.Fatalf("title expected to be foo but is %s", title)
    }
}
```

## Testing for errors

You can also mock the client to return an error for a given query by using the `Errors` function.

```go
// main_test.go
func TestGetPostTitle_error(t *testing.T) {
    client, mock, ensure := NewMock()
    defer ensure(t)

    mock.Post.Expect(
        client.Post.FindUnique(
            db.Post.ID.Equals("123"),
        ),
    ).Errors(db.ErrNotFound)

    _, err := GetPostTitle(context.Background(), client, "123")
    if !errors.Is(err, ErrNotFound) {
        t.Fatalf("error expected to return ErrNotFound but is %s", err)
    }
}
```

## Next steps

Learn how to build up [dynamic queries](17-dynamic-queries.md).
