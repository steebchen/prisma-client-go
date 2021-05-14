# Introduction

The Prisma Go client aims to be fully type-safe wherever possible, even for complex queries. That's why it uses a
functional syntax to make sure every argument you provide matches the type in the database. If you change something in
the database, and the change is incompatible with the code, you will get a compile time error.

The examples use the following prisma schema:

```prisma
model Post {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    updatedAt DateTime @updatedAt
    published Boolean
    title     String
    content   String?

    comments Comment[]
}

model Comment {
    id        String   @default(cuid()) @id
    createdAt DateTime @default(now())
    content   String

    post   Post @relation(fields: [postID], references: [id])
    postID String
}
```

## The syntax concept

The generated Go client uses a specific and consistent functional API syntax.

To fetch all posts, you can do:

```go
posts, err := client.Post.FindMany().Exec(ctx)
```

Usually you want to query for some ID or other field. The Go client exposes a completely type-safe query builder in the
form of `<model>.<field>.<method>`:

```go
posts, err := client.Post.FindMany(
    // package.
    // model.
    //      field.
    //            method.
    //                  (value)
    db.Post.Title.Equals("hi"),
).Exec(ctx)
```

This schema is consistently used. You can usually just type `db.` and then see what models, fields and actions your editor auto completion will suggest. Depending on the type or the shape of a model or a field, there may be different actions available. If a field is optional, you will also get additional
methods such as IsNull() and *Optional variations to query for SQL NULLs:

```go
posts, err := client.Post.FindMany(
    db.Post.Title.Equals("hi"),
    db.Post.Title.Contains("hi"),
    db.Post.Content.IsNull(),
    db.Post.Desc.Contains(variable),
    db.Post.Desc.ContainsOptional(pointerVariable),
).Exec(ctx)
```

## Next steps

We'll explore how you can query for data in the [next article](02-find.md).
