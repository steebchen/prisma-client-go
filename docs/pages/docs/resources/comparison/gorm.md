# Prisma Client Go vs gorm

This article assumes you have basic knowledge of ORMs. If not, start with
our [Next Gen ORM](../../resources/next-gen-orm.md) article first.

## GORM – the fantastic database library, or merely a basic query builder?

### About Gorm

[Gorm](https://gorm.io) describes itself as a fantastic ORM and has been around for a long time. This article analyzes this and compares Gorm with Prisma Client Go.

### Basics

Let's look into the details: setup, how data is queried or modified, features and how migrations work.

#### Key differences in the query builder

Gorm is mostly a query builder, which makes it lightweight, but also means that big parts of your query are not fully typesafe, and makes you learn extra syntax. Prisma has a bigger generated client, but gives you all the flexibility and query capabilities while ensuring everything is still typesafe.

#### Querying for data

A simple query for a table User will look as follows:

```go
// Gorm
user, err := .User()
```

```go
// Prisma
user, err := client.User.FindUnique(
  db.User.Email.Equals("gopher@example.com"),
)
```

As you can see, the Prisma query is a little bit more verbose, but it's not just more clear at what it does, it also uses more of the generated type-safe methods, versus requiring you to combine operations using string as Gorm does.

#### Writing data

Writing data works as follows:

```go
// Gorm
user, err := .User()
```

```go
// Prisma
user, err := client.User.FindUnique(
  db.User.Email.Equals("gopher@example.com"),
)
```

#### Features

#### Migrations

Gorm provides simple migrations, which means that upon running your app, it can create the tables according to your defined types.

However, this quickly falls apart when you want to modify any existing columns, especially once you have production data. To properly handle this, you would need to introduce another third party tool, configure it, figure out how it works, and keep your Gorm model as well as your migrations in sync. There needs to be a better way!

Prisma, on the other hand, has built-in migrations for multiple databases. For prototyping, you can easily sync your wanted database layout to your database, but for your production database, there is a fully fledged migration tool.

To use it, you can create migrations, optionally adapt them if you have some custom additions, and then just run `migrate deploy` on your production database. The migration is wrapped in transaction, it remembers which migrations already ran, and if there are any conflicts in your production database, you will get a clear error, after which you can fix the migration and re-deploy.

It gets better, though: if you already have an existing database, you can just use the `db pull` command to create a schema from your existing production database, and you will have a fully typesafe database client generated for you, regardless of how big your schema is or how many tables you have.
