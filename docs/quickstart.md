# Quickstart

## Setup

1) Initialise a new Go project

    If you don't have a Go project yet, initialise one using Go modules:

    ```shell script
    mkdir demo && cd demo
    go mod init demo
    ```

2) Get Prisma Client Go

    Install the Go module in your project:

    ```shell script
    go get github.com/prisma/prisma-client-go
    ```

3) Prepare your database schema in a `schema.prisma` file. For example, a simple schema with a sqlite database and
    Prisma Client Go as a generator with two models would look like this:

    ```prisma
    datasource db {
        // could be postgresql or mysql
        provider = "sqlite"
        url      = "file:dev.db"
    }

    generator db {
        provider = "go run github.com/prisma/prisma-client-go"
    }

    model Post {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        updatedAt DateTime @updatedAt
        title     String
        published Boolean
        desc      String?
    }
    ```

    To get this up and running in your database, we use the Prisma migration
    tool [`migrate`](https://www.prisma.io/docs/concepts/components/prisma-migrate) (Note: this tool is experimental) to create and migrate our
    database:

     ```shell script
    # sync the database with your schema
    go run github.com/prisma/prisma-client-go db push
    ```

4) Generate the Prisma Client Go client in your project

    ```shell script
    go run github.com/prisma/prisma-client-go generate
    ```

    If you make changes to your prisma schema, you need to run this command again.

## Usage

Create a file `main.go`:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"

    "demo/db"
)

func main() {
    if err := run(); err != nil {
        panic(err)
    }
}

func run() error {
    client := db.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        return err
    }

    defer func() {
        if err := client.Prisma.Disconnect(); err != nil {
            panic(err)
        }
    }()

    ctx := context.Background()

    // create a post
    createdPost, err := client.Post.CreateOne(
        db.Post.Title.Set("Hi from Prisma!"),
        db.Post.Published.Set(true),
        db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
    ).Exec(ctx)
    if err != nil {
        return err
    }

    result, _ := json.MarshalIndent(createdPost, "", "  ")
    fmt.Printf("created post: %s\n", result)

    // find a single post
    post, err := client.Post.FindUnique(
        db.Post.ID.Equals(createdPost.ID),
    ).Exec(ctx)
    if err != nil {
        return err
    }

    result, _ = json.MarshalIndent(post, "", "  ")
    fmt.Printf("post: %s\n", result)

    // for optional/nullable values, you need to check the function and create two return values
    // `desc` is a string, and `ok` is a bool whether the record is null or not. If it's null,
    // `ok` is false, and `desc` will default to Go's default values; in this case an empty string (""). Otherwise,
    // `ok` is true and `desc` will be "my description".
    desc, ok := post.Desc()
    if !ok {
        return fmt.Errorf("post's description is null")
    }

    fmt.Printf("The posts's description is: %s\n", desc)

    return nil
}
```

and run it:

```shell script
go run .
```

```
❯ go run .
created post: {
  "id": "ckfnrp7ec0000oh9kygil9s94",
  "createdAt": "2020-09-29T09:37:44.628Z",
  "updatedAt": "2020-09-29T09:37:44.628Z",
  "title": "Hi from Prisma!",
  "published": true,
  "desc": "Prisma is a database toolkit and makes databases easy."
}
post: {
  "id": "ckfnrp7ec0000oh9kygil9s94",
  "createdAt": "2020-09-29T09:37:44.628Z",
  "updatedAt": "2020-09-29T09:37:44.628Z",
  "title": "Hi from Prisma!",
  "published": true,
  "desc": "Prisma is a database toolkit and makes databases easy."
}
The posts's title is: Prisma is a database toolkit and makes databases easy.
```

### Next steps

We just scratched the surface of what you can do. Read our [advanced tutorial](advanced.md) to learn about more
complex queries and how you can query for relations.
