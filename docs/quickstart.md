# Quickstart

## Set up Prisma

1) Create a new project with Go modules

    Skip this step if you're using an existing project.

    ```shell script
    go mod init github.com/your/repo
    ```

2) Get Prisma Client Go

    Prisma client Go is decoupled from Prisma in a way that you can use it without manually instally the Prisma CLI. Instead, it is shipped with the Go module and downloaded for you.

    ```shell script
    go get github.com/prisma/prisma-client-go
    ```

3) [Prepare your Prisma schema](https://www.prisma.io/docs/reference/tools-and-interfaces/prisma-schema/prisma-schema-file) in a `schema.prisma` file. For example, a simple schema with a sqlite database and Prisma Client Go as a generator with two models would look like this:

    ```prisma
    datasource db {
        provider = "sqlite"
        url      = "file:dev.db"
    }

    generator db {
        provider = "go run github.com/prisma/prisma-client-go"
    }

    model User {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        email     String   @unique
        name      String?
        age       Int?

        posts     Post[]
    }

    model Post {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        updatedAt DateTime @updatedAt
        published Boolean
        title     String
        content   String?

        author   User @relation(fields: [authorID], references: [id])
        authorID String
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

    Prisma Client Go is now generated into the file path you specified in the "output" option which is `"./db/db_gen.go"` in this case.
    If you make changes to your prisma schema, you need to run this command again.

## Usage

Once you generated the Prisma Client Go client and set up a datasource with Prisma, you're good to go!

We recommend generating the client into a package called `db` (see step 3 above) at `./db/db_gen.go`, but you can adapt these settings to anything you want.

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
    "github.com/your/repo/db"
)

func main() {
    client := db.NewClient()
    err := client.Connect()
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

    // create a user
    createdUser, err := client.User.CreateOne(
        db.User.Email.Set("john.doe@example.com"),
        db.User.Name.Set("John Doe"),

        // ID is optional, which is why it's specified last. if you don't set it
        // an ID is auto generated for you
        db.User.ID.Set("123"),
    ).Exec(ctx)

    log.Printf("created user: %+v", createdUser)

    // find a single user
    user, err := client.User.FindOne(
        db.User.Email.Equals("john.doe@example.com"),
    ).Exec(ctx)
    if err != nil {
        panic(err)
    }

    log.Printf("user: %+v", user)

    // for optional/nullable values, you need to check the function and create two return values
    // `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
    // `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
    // `ok` is true and `name` will be "John Doe".
    name, ok := user.Name()

    if !ok {
        log.Printf("user's name is null")
        return
    }

    log.Printf("The users's name is: %s", name)
}
```
