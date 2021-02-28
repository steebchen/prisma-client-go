# JSON

The examples use the following prisma schema:

```prisma
model Post {
    id    String   @default(cuid()) @id
    title String
    info  Json
}
```

## How JSON works in the Go client

JSON works by using Go's `json.RawMessage` data structure. It's a `[]byte` behind the scenes, which means you can use the API you already know to work with unstructured json data in the Go client.

You can work with []bytes directly, but usually you marshal this data from an existing struct or unmarshal the data to a given struct variable.

## Write JSON data

```go
type PostInfo struct {
    Content string `json:"content"`
}

postInfo := &PostInfo{
    Content: "hi",
}
infoBytes, err := json.Marshal(postInfo)
if err != nil {
    panic(err)
}

_, err = client.Post.CreateOne(
    Post.Title.Set("what up"),
    Post.Info.Set(infoBytes),
    Post.ID.Set("123"),
).Exec(ctx)
if err != nil {
    panic(err)
}
```

## Read JSON data

```go
post, err := client.Post.FindUnique(
    Post.ID.Equals("123"),
).Exec(ctx)
if err != nil {
    panic(err)
}

// post.Info is of type json.RawMessage, so this will contain binary data such as [123 34 97 116 116 ...]
// however, if we format it with %s, we can convert the contents to a string to see what's inside:
log.Printf("post info: %s", post.Info)

// to unmarshal this information into a specific struct, we make use of Go's usual handling of json data:

type PostInfo struct {
    Content string `json:"content"`
}

var info PostInfo
if err := json.Unmarshal(post.Info, &info); err != nil {
    panic(err)
}
log.Printf("post info: %+v", info)
```

## Next steps

Check out how to [dynamically ignore fields](15-if-present-methods.md).
