# Photon Go

Photon Go is an auto-generated database client, which is fully typesafe, reduces boilerplate code and replaces traditional ORMs.

Photon Go is a part of the [Prisma Framework](https://github.com/prisma/prisma2) and depends on it.

*NOTE*: Currently, Photon Go is still under heavy development and in an experimentation phase, where we build specifications for the syntax and make major decisions on how the Go client will look and work. Photon Go is unstable and there is no ETA for general availability yet.

## Setup

The Photon Go setup has a few caveats due to Prisma Framework being in an alpha state at this time.
You can already use Photon Go, but it means you have to take a few extra steps. This will not be required anymore in future releases.

1) Install the Prisma CLI (requires Node.JS/npm)
    ```
    npm i -g prisma2
    ```

2) Get Photon Go
    ```
    go get github.com/prisma/photongo
    ```

3) [Prepare your Prisma schema](https://github.com/prisma/prisma2/blob/master/docs/prisma-schema-file.md) in a `schema.prisma` file, and specify `photongo` as a provider:

    ```
    generator photon {
      provider = "photongo"
      output = "./photon/photon_gen.go"
      package = "photon"
    }
    ```

4) Generate the Photon Go client in your project
    ```
    prisma2 generate
    ```

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

### Reference

The photon client provides the methods FindOne, FindMany, and CreateOne. Just with these 3 methods, you can query for anything, and optionally update or delete for the queried records.

Additionally, Photon Go provids a fully fluent and type-safe query API, which always follows the schema `photon.<Model>.<Field>.<Action>`, e.g. `photon.User.Name.Equals("John")`.

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
  User.ID.Set("id"),
  User.Email.Set("email"),
  User.Username.Set("username"),

  // optional fields
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
).Update(
  User.Username.Set("new-username"),
  User.Name.Set("New Name"),
).Exec(ctx)
```

#### Quering for relations
