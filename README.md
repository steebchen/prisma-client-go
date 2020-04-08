# Prisma Client Go

Prisma Client Go is an auto-generated database client, which is fully typesafe, reduces boilerplate code and replaces traditional ORMs.

Prisma Client Go is a part of [Prisma](https://github.com/prisma/prisma2) and depends on it.

**NOTE: Currently, the Prisma Go Client is still under heavy development and in an experimentation phase. It's a prototype, and there can and will be breaking changes. Prisma Client Go is unstable and there is no ETA for general availability yet. The current API is not final and may change.**

We recommend to read the [current caveats](#caveats).

## Setup

1) Get Prisma Client Go
    ```shell script
    go get github.com/prisma/prisma-client-go
    ```

2) [Prepare your Prisma schema](https://www.prisma.io/docs/reference/tools-and-interfaces/prisma-schema/prisma-schema-file) in a `schema.prisma` file. For example, a simple schema with a sqlite database and Prisma Client Go as a generator with two models would look like this:

    ```prisma
    datasource db {
      provider = "sqlite"
      url      = "file:dev.db"
    }

    generator db {
      provider = "go run github.com/prisma/prisma-client-go"
      output = "./db/db_gen.go"
      package = "db"
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

    To get this up and running in your database, we use the Prisma migration tool [`migrate`](https://github.com/prisma/migrate) to create and migrate our database:

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

Note: Some errors may get displayed, but you can ignore them. Prisma Studio is currently not working. As an alternative, you can install the [Prisma CLI](https://github.com/prisma/prisma2#getting-started).

For more information and instructions on how to deploy your app, please check the [deploy instructions](#deploy).

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

### Example

```go
import (
  "context"
  "fmt"
  "github.com/your/repo/db"
)

func main() {
  client := db.NewClient()
  err := client.Connect()
  check(err)

  defer func() {
    err := client.Disconnect()
    check(err)
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

  fmt.Printf("created user: %+v\n", createdUser)

  // find a single user
  user, err := client.User.FindOne(
    db.User.Email.Equals("john@example.com"),
  ).Exec(ctx)
  check(err)

  fmt.Printf("user: %+v\n", user)

  // for optional/nullable values, you need to check the function and create two return values
  // `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
  // `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
  // `ok` is true and `name` will be "John Doe".
  name, ok := user.Name()

  if !ok {
    fmt.Printf("user's name is null\n")
    return
  }

  fmt.Printf("The users's name is: %s\n", name)
}
```

### Deploy

Deploying a Prisma Client Go adds a few more steps, because it depends on the Prisma query engine, which is a binary we automatically download in your project path. Depending on where you deploy your code to, you might need to follow some extra steps.

#### Set up go generate

While this step is not required, we recommend to use [`go generate`](https://blog.golang.org/generate) to simplify generating the Prisma Client Go client. To do so, you can just put the following line into a go file, and then run go generate so `go run github.com/prisma/prisma-client-go generate` and any other generate commands you run will get executed.

Put this line into a go file in your project, usually in `main.go`:

```go
//go:generate go run github.com/prisma/prisma-client-go generate

func main() {
  // ...
}
```

Now, run `go generate`:

```shell script
go generate
```

Your Prisma Client Go code is now generated.

#### Traditionally deploy to a server

Usually, you would deploy your Go app by running `go build .`, which generates a binary, and then deploy that binary anywhere you want. However, since Prisma Client Go depends on the Prisma query engine, you also need to deploy the query engine binary `query-engine-*` files.

If you use different development environments, e.g. a Mac to develop, and Debian on your server, you need to specify these two binaries in the schema.prisma file so that you can then also upload the binary suitable for your deploy environment.

```prisma
generator db {
  provider = "go run github.com/prisma/prisma-client-go"
  binaryTargets = ["native", "debian-openssl-1.1.x"]
}
```

You can find all binary targets [in our specs repository](https://github.com/prisma/specs/tree/master/binaries#binary-builds).

#### Using docker

When deploying with docker, the setup is super easy. Build your dockerfile as usual, run `go generate ./...` (see [setting up go generate](#set-up-go-generate)), and you're good to go!

We also recommend using [Go modules](https://blog.golang.org/using-go-modules), which is recommended when using Go >=1.13.

Your dockerfile could look like this. It uses Go modules, layered caching for fast docker builds and multiple stages for lightweight images (usually a few megabytes).

```dockerfile
FROM golang:1.13 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

COPY . ./

# generate the Prisma Client Go client
RUN go run github.com/prisma/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /main .

CMD ["/main"]
```

### Reference

The db client provides the methods FindOne, FindMany, and CreateOne. Just with these 3 methods, you can query for anything, and optionally update or delete for the queried records.

Additionally, Prisma Client Go provides a fully fluent and type-safe query API, which always follows the schema `db.<Model>.<Field>.<Action>`, e.g. `db.User.Name.Equals("John")`.

#### Reading data

##### Find many records

```go
users, err := client.User.FindMany(
  photon.User.Name.Equals("hi"),
).Exec(ctx)
```

If no records are found, this returns an empty array without returning an error (like usual SQL queries).

##### Find one record

```go
user, err := client.User.FindOne(
  db.User.ID.Equals("123"),
).Exec(ctx)

if err == db.ErrNotFound {
  log.Printf("no record with id 123")
}
```

This returns an error of type `ErrNotFound` (exported in the `db` package) if there was no such record.

##### Query API

Depending on the data types of your fields, you will automatically be able to query for respective operations. For example, for integer or float fields you might want to query for a field which is less than or greater than some number.

```go
user, err := client.User.FindOne(
  // query for names containing the string "Jo"
  db.User.Name.Contains("Jo"),
).Exec(ctx)
```

Other possible queries are:

```go
// query for people who are named "John"
db.User.Name.Contains("John"),
// query for names containing the string "oh"
db.User.Name.Contains("oh"),
// query for names starting with "Jo"
db.User.Name.HasPrefix("Jo"),
// query for names ending with "Jo"
db.User.Name.HasSuffix("hn"),
// query for all users which are younger than or exactly 18
db.User.Age.LTE(18),
// query for all users which are younger than 18
db.User.Age.LT(18),
// query for all users which are older than or exactly 18
db.User.Age.GT(18),
// query for all users which are older than 18
db.User.Age.GTE(18),
// query for all users which were created in the last 6 hours
db.User.CreatedAt.After(time.Now().Add(-6 * time.Hour)),
// query for all users which were created until yesterday
db.User.CreatedAt.Before(time.Now().Truncate(24 * time.Hour)),
```

All of these queries are fully type-safe and independent of the underlying database.

#### Writing data

##### Create a record

```go
created, err := client.User.CreateOne(
  // required fields
  User.Email.Set("email"),
  User.Username.Set("username"),

  // optional fields
  User.ID.Set("id"),
  User.Name.Set("name"),
  User.Stuff.Set("stuff"),
).Exec(ctx)
```

##### Update a record

To update a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Update()`.

```go
updated, err := client.User.FindOne(
  User.Email.Equals("john@example.com"),
).Update(
  User.Username.Set("new-username"),
  User.Name.Set("New Name"),
).Exec(ctx)
```

##### Delete a record

To delete a record, just query for a field using FindOne or FindMany, and then just chain it by invoking `.Delete()`.

```go
updated, err := client.User.FindOne(
  User.Email.Equals("john@example.com"),
).Delete().Exec(ctx)
```

#### Querying for relations

*TBD*

## Caveats

Prisma Client Go is experimental and comes with some caveats. We plan to eliminate all of these in the future.

- We recommend to use Go 1.13 or higher, everything else is untested.
- Expect breaking changes in minor versions in 0.x.x releases.
- Not tested on Windows.
- DO NOT USE IN PRODUCTION!
