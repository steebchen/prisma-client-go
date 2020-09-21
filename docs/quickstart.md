# Quickstart

## Setup

1) Init go project

    If you don't have a go project yet, initialise one using go modules:

    ```shell script
    mkdir demo && cd demo
    go mod init demo
    ```

2) Get Prisma Client Go

    Install the go module in your project:

    ```shell script
    go get github.com/prisma/prisma-client-go
    ```

3) Prepare your database schema in a `schema.prisma` file. For example, a simple schema with a sqlite database and Prisma Client Go as a generator with two models would look like this:

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
        published Boolean
        title     String
        content   String?
    }
    ```

    To get this up and running in your database, we use the Prisma migration tool [`migrate`](https://github.com/prisma/migrate) (Note: this tool is experimental) to create and migrate our database:

    ```shell script
    # initialize the first migration
    go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name "init"
    # apply the migration
    go run github.com/prisma/prisma-client-go migrate up --experimental
    ```

4) Generate the Prisma Client Go client in your project

    ```shell script
    go run github.com/prisma/prisma-client-go generate
    ```

    If you make changes to your prisma schema, you need to run this command again.

## Usage

Once you generated the Prisma Client Go client and set up a datasource with Prisma, you're good to go!

### Create the client and connect to the prisma engine

```go
client := db.NewClient()
err := client.Connect()
if err != nil {
    handle(err)
}

defer func() {
    err := client.Disconnect()
    if err != nil {
        panic(fmt.Errorf("could not disconnect %w", err))
    }
}()
```

### Full example

```go
package main

import (
    "context"
    "log"

    "demo/db"
)

func main() {
    client := db.NewClient()
    err := client.Connect()``
    if err != nil {
        panic(err)
    }

    defer func() {
        err := client.Disconnect()
        if err != nil {
            panic(err)
        }
    }()

    ctx := context.Background()

    // create a post
    createdPost, err := client.Post.CreateOne(
        db.Post.Title.Set("Hi from Prisma!"),
        db.Post.Published.Set(true),
        db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
        // ID is optional since it's auto generated, which is why it's specified last.
        db.Post.ID.Set("123"),
    ).Exec(ctx)
    if err != nil {
        panic(err)
    }

    log.Printf("created post: %+v", createdPost)

    // find a single post
    post, err := client.Post.FindOne(
        db.Post.Email.Equals("john.doe@example.com"),
    ).Exec(ctx)
    if err != nil {
        panic(err)
    }

    log.Printf("post: %+v", post)

    // for optional/nullable values, you need to check the function and create two return values
    // `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
    // `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
    // `ok` is true and `desc` will be "my description".
    name, ok := post.Desc()

    if !ok {
        log.Printf("post's name is null")
        return
    }

    log.Printf("The posts's name is: %s", name)
}
```

### Next steps

We just scratched the surface of what you can do. Read our [advanced tutorial](./advanced.md) to
learn about more complex queries and how you can query for relations.
