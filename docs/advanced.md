## Advanced usage

In the [quickstart](./quickstart.md), we have created a simple user model and ran a few queries.
However, Prisma and the Go client are designed to work with relations between models.

We already created a post model, such as for a blog. Let's assume we want to add a user model, and connect these models
in a way so we can rely on SQL's foreign keys and the Go client's ability to work with relations.

So let's introduce a user model:

```prisma
model User {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    email     String   @unique
    name      String?
    age       Int?

    posts     Post[]
}
```

We will also need to add a relation from to the post model in order to make a 1:n relation between those models:

```prisma
model Post {
    // ...

    author   User @relation(fields: [authorID], references: [id])
    authorID String
}
```

<details>
    <summary>Expand to show full schema.prisma</summary>

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

        author   User @relation(fields: [authorID], references: [id])
        authorID String
    }

    model User {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        email     String   @unique
        name      String?
        age       Int?

        posts     Post[]
    }
    ```
</details>

Whenever you make changes to your model, migrate your database and re-generate your prisma code:

```shell script
# apply migrations
go run github.com/prisma/prisma-client-go migrate save --experimental --name "add user model"
go run github.com/prisma/prisma-client-go migrate up --experimental
# generate
go run github.com/prisma/prisma-client-go generate
```

In order to create a post, we first need to create a user, and then reference that user when creating a post.

```go
// create a user first
user, err := client.User.CreateOne(
    db.User.Email.Set("john.doe@example.com"),
    db.User.Name.Set("John Doe"),
    db.Post.Desc.Set("Hi there."),
).Exec(ctx)

// create a post and set the author
_, err := client.Post.CreateOne(
    db.Post.Title.Set("My new post"),
    db.Post.Published.Set(true),
    db.Post.Desc.Set("Hi there."),
    db.Post.Author.Link(
        db.User.ID.Equals(user.ID),
    ),
).Exec(ctx)
```

Now that a post and a user row are created, you can query for them as follows:

```go
// return all posts which belong to author with id 123
_, err := client.Post.FindMany(
    db.Post.Author.Where(
        db.User.ID.Equals("123"),
    ),
).Exec(ctx)

// return all users which have a post containing the word "title"
_, err := client.User.FindMany(
    db.User.Posts.Some(
        db.Post.Title.Contains("post"),
    ),
).Exec(ctx)
```

Prisma also allows you to fetch multiple things at once. Instead of doing complicated joins, you can fetch a user and
a few of their posts in just a few lines and fully typesafe:

```go
// return all posts which belong to author with id 123
_, err := client.Post.FindMany(
    db.Post.Author.Where(
        db.User.ID.Equals("123"),
    ),
).Exec(ctx)

// return all users which have a post containing the word "title"
_, err := client.User.FindMany(
    db.User.Posts.Some(
        db.Post.Title.Contains("post"),
    ),
).Exec(ctx)
```

## API reference

To explore all query capabilities, check out the [API reference](./reference).
