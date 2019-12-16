# Photon Go

Photon Go is an auto-generated database client, which is fully typesafe, reduces boilerplate code and replaces traditional ORMs.

Photon Go is a part of [Prisma](https://github.com/prisma/prisma2) and depends on it.

**NOTE: Currently, Photon Go is still under heavy development and in an experimentation phase. It's a prototype, and there can and will be breaking changes. Photon Go is unstable and there is no ETA for general availability yet.**

We recommend to read the [current caveats](#caveats).

## Setup

1) Get Photon Go
    ```shell script
    go get github.com/prisma/photongo
    ```

2) [Prepare your Prisma schema](https://github.com/prisma/prisma2/blob/master/docs/prisma-schema-file.md) in a `schema.prisma` file. For example, a simple schema with a sqlite database and Photon Go as a generator with two models would look like this:

    ```prisma
    datasource db {
      provider = "sqlite"
      url      = "file:dev.db"
    }

    generator photon {
      provider = "photongo"
      output = "./photon/photon_gen.go"
      package = "photon"
    }

    model User {
      id        String   @default(cuid()) @id @unique
      createdAt DateTime @default(now())
      email     String   @unique
      name      String?
      age       Int?
      posts     Post[]
    }

    model Post {
      id        String   @default(cuid()) @id @unique
      createdAt DateTime @default(now())
      updatedAt DateTime @updatedAt
      published Boolean
      title     String
      content   String?
      author    User
    }
    ```

    To get this up and running in your database, we use the Prisma migration tool [`lift`](https://github.com/prisma/lift) to create and migrate our database:

    ```shell script
    # initialize the first migration
    go run github.com/prisma/photongo lift save --create-db --name "init"
    # apply the migration
    go run github.com/prisma/photongo lift up
    ```

4) Generate the Photon Go client in your project

    ```shell script
    go run github.com/prisma/photongo generate
    ```

    Photon go is now generated into the file path you specified in the "output" option which is `"./photon/photon_gen.go"` in this case.

For development, you can also use the dev command for continuous generation. It will also automatically handle migrations locally whenever you change your schema.

```shell script
go run github.com/prisma/photongo dev
```

Note: Some errors may get displayed, but you can ignore them. Prisma Studio is currently not working. As an alternative, you can install the [Prisma CLI](https://github.com/prisma/prisma2#getting-started).

For more information and instructions on how to deploy your app, please check the [deploy instructions](#deploy).

## Usage

Once you generated the Photon Go client and set up a datasource with Prisma, you're good to go!

We recommend generating the client into a package called `photon` (see step 3 above) at `./photon/photon_gen.go`, but you can adapt these settings to anything you want.

### Create the client and connect to the prisma engine

```go
client := photon.NewClient()
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
  "github.com/your/repo/photon"
)

func main() {
  client := photon.NewClient()
  err := client.Connect()
  check(err)

  defer func() {
    err := client.Disconnect()
    check(err)
  }()

  ctx := context.Background()

  // create a user
  createdUser, err := client.User.CreateOne(
    photon.User.Email.Set("john.doe@example.com"),
    photon.User.Name.Set("John Doe"),

    // ID is optional, which is why it's specified last. if you don't set it
    // an ID is auto generated for you
    photon.User.ID.Set("123"),
  ).Exec(ctx)

  fmt.Printf("created user: %+v\n", createdUser)

  // find a single user
  user, err := client.User.FindOne(
    photon.User.Email.Equals("john@example.com"),
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

Deploying a Photon Go adds a few more steps, because it depends on the Prisma query engine, which is a binary we automatically download in your project path. Depending on where you deploy your code to, you might need to follow some extra steps.

#### Set up go generate

While this step is not required, we recommend to use [`go generate`](https://blog.golang.org/generate) to simplify generating the Photon Go client. To do so, you can just put the following line into a go file, and then run go generate so `go run github.com/prisma/photongo generate` and any other generate commands you run will get executed.

Put this line into a go file in your project, usually in `main.go`:

```go
//go:generate go run github.com/prisma/photongo generate

func main() {
  // ...
}
```

Now, run `go generate`:

```shell script
go generate
```

Your Photon Go code is now generated.

#### Traditionally deploy to a server

Usually, you would deploy your Go app by running `go build .`, which generates a binary, and then deploy that binary anywhere you want. However, since Photon Go depends on the Prisma query engine, you also need to deploy the query engine binary `query-engine-*` files.

If you use different development environments, e.g. a Mac to develop, and Debian on your server, you need to specify these two binaries in the schema.prisma file so that you can then also upload the binary suitable for your deploy environment.

```prisma
generator photon {
  provider = "photongo"
  binaryTargets = ["native", "debian-openssl-1.1.x"]
}
```

You can find all binary targets [in our specs repository](https://github.com/prisma/specs/tree/master/binaries#binary-builds).

#### Using docker

When deploying with docker, the setup is super easy. Build your dockerfile as usual, run `go generate` (see [setting up go generate](#set-up-go-generate)), and you're good to go!

We also recommend using [Go modules](https://blog.golang.org/using-go-modules), which is recommended when using Go >=1.13.

Your dockerfile could look like this. It uses Go modules, layered caching for fast docker builds and multiple stages for lightweight images (usually a few megabytes).

```dockerfile
FROM golang:1.13 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# temporarily needed to enforce installing the photongo binary
RUN go install github.com/prisma/photongo

COPY . ./

# generate the Photon Go client
RUN go run github.com/prisma/photongo generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /main .

CMD ["/main"]
```

### Reference

The photon client provides the methods FindOne, FindMany, and CreateOne. Just with these 3 methods, you can query for anything, and optionally update or delete for the queried records.

Additionally, Photon Go provides a fully fluent and type-safe query API, which always follows the schema `photon.<Model>.<Field>.<Action>`, e.g. `photon.User.Name.Equals("John")`.

#### Reading data

##### Find many records

```go
user, err := client.User.FindOne(
  photon.User.Name.Equals("hi"),
).Exec(ctx)
```

If no records are found, this returns an empty array without returning an error (like usual SQL queries).

##### Find one record

```go
user, err := client.User.FindOne(
  photon.User.ID.Equals("123"),
).Exec(ctx)

if err == photon.ErrNotFound {
  log.Printf("no record with id 123")
}
```

This returns an error of type `ErrNotFound` (exported in the `photon` package) if there was no such record.

##### Query API

Depending on the data types of your fields, you will automatically be able to query for respective operations. For example, for integer or float fields you might want to query for a field which is less than or greater than some number.

```go
user, err := client.User.FindOne(
  // query for names containing the string "Jo"
  photon.User.Name.Contains("Jo"),
).Exec(ctx)
```

Other possible queries are:

```go
// query for people who are named "John"
photon.User.Name.Contains("John"),
// query for names containing the string "oh"
photon.User.Name.Contains("oh"),
// query for names starting with "Jo"
photon.User.Name.HasPrefix("Jo"),
// query for names ending with "Jo"
photon.User.Name.HasSuffix("hn"),
// query for all users which are younger than or exactly 18
photon.User.Age.LTE(18),
// query for all users which are younger than 18
photon.User.Age.LT(18),
// query for all users which are older than or exactly 18
photon.User.Age.GT(18),
// query for all users which are older than 18
photon.User.Age.GTE(18),
// query for all users which were created in the last 6 hours
photon.User.CreatedAt.After(time.Now().Add(-6 * time.Hour)),
// query for all users which were created until yesterday
photon.User.CreatedAt.Before(time.Now().Truncate(24 * time.Hour)),
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

Photon Go is experimental and comes with some caveats. We plan to eliminate all of these in the future.

- We recommend to use Go 1.13 or higher, everything else is untested.
- You need to use an absolute path for the generator. Tracked in https://github.com/prisma/prisma2/issues/934.
- Multiple projects using Photon Go can conflict because a given Go binary just defaults to the last installed version (https://github.com/golang/go/issues/27653). We recommend to use the same version in all your projects, even if they are unrelated. Waiting for https://github.com/prisma/prisma2/issues/1101.
- Expect breaking changes in minor versions in 0.x.x releases.
- No Windows support.
- DO NOT USE IN PRODUCTION!
