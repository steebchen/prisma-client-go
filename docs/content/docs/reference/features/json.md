# JSON

The examples use the following prisma schema:

```prisma
model Log {
  id      String   @id @default(cuid())
  date    DateTime @default(now())
  message String
  meta    Json
}
```

## How JSON works in the Go client

JSON works by using Go's `json.RawMessage` data structure. It's a `[]byte` behind the scenes, which means you can use
the API you already know to work with unstructured json data in the Go client.

You can work with []bytes directly, but usually you marshal this data from an existing struct or unmarshal the data to a
given struct variable.

## Write JSON data

```go
type LogInfo struct {
  Service string `json:"service"`
}

logInfo := &LogInfo{
  Service: "deployment/api",
}
infoBytes, err := json.Marshal(logInfo)
if err != nil {
  panic(err)
}

_, err = client.Log.CreateOne(
  db.Log.Message.Set("/api/graphql: status code 400"),
  db.Log.Info.Set(infoBytes),
  db.Log.ID.Set("123"),
).Exec(ctx)
if err != nil {
  panic(err)
}
```

## Read JSON data

```go
log, err := client.Log.FindUnique(
  db.Log.ID.Equals("123"),
).Exec(ctx)
if err != nil {
  panic(err)
}

// log.Info is of type json.RawMessage, so this will contain binary data such as [123 34 97 116 116 ...]
// however, if we format it with %s, we can convert the contents to a string to see what's inside:
log.Printf("log info: %s", log.Info)

// to unmarshal this information into a specific struct, we make use of Go's usual handling of json data:

type LogInfo struct {
  Service string `json:"service"`
}

var info LogInfo
if err := json.Unmarshal(log.Info, &info); err != nil {
  panic(err)
}
log.Printf("log info: %+v", info)
```

## Query JSON

You can filter JSON fields by using a combination of `Path` and a JSON query. Note that the syntax differs between
databases.

```go
actual, err := client.User.FindFirst(
  User.Meta.Path([]string{"service"}),
  User.Meta.StringContains("api"),
).Exec(ctx)
```

```go
actual, err := client.User.FindFirst(
  User.Meta.Path([]string{"service"}),
  // Note that Equals accepts JSON, so strings need to be surrounded with quotes
  User.Meta.Equals(JSON(`"deployment/api"`)),
).Exec(ctx)
```

For more information about all json filters and more example queries, check out
the [Prisma JSON filters documentation](https://www.prisma.io/docs/concepts/components/prisma-client/working-with-fields/working-with-json-fields).
