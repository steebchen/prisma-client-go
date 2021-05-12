<br />

<div align="center">
    <h1>Prisma Client Go</h1>
    <p><h3 align="center">Typesafe database access for Go</h3></p>
    <a href="./docs/quickstart.md">Quickstart</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://www.prisma.io/">Website</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="./docs">Docs</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="./docs/reference">API reference</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://www.prisma.io/blog">Blog</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://slack.prisma.io/">Slack</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://twitter.com/prisma">Twitter</a>
</div>

<hr>

Prisma Client Go is an **auto-generated query builder** that enables **type-safe** database access and **reduces boilerplate**. You can use it as an alternative to traditional ORMs such as gorm, xorm, sqlboiler and most database-specific tools.

It is part of the [Prisma](https://www.prisma.io/) ecosystem. Prisma provides database tools for data access, declarative data modeling, schema migrations and visual data management.

_NOTE_: Prisma Client Go is currently offered under our [early access program](https://www.prisma.io/docs/about/releases#product-maturity-levels). There will be documented breaking changes with new [releases](https://github.com/prisma/prisma-client-go/releases).

## Getting started

To get started, [**read our quickstart tutorial**](./docs/quickstart.md) to add Prisma to your project in just a few minutes.

You also might want to read [deployment tips](./docs/deploy.md) and the [full API reference](./docs/reference).

## Notes

The go client works slightly different than the normal Prisma tooling. When you're using the go client, whenever you see Prisma CLI commands such as `prisma ...`, you should always write `go run github.com/prisma/prisma-client-go ...` instead.

## Contributing

### Running Tests

```shell
# requires docker to be installed
go run ./test/setup/init setup # sets up docker containers for integration testing
go generate ./...
go test ./... -v

# to teardown docker containers:
go run ./test/setup/init teardown
```

### Writing Commit Messages

We use [conventional commits](https://www.conventionalcommits.org) (also known as semantic commits) to ensure consistent and descriptive commit messages.

## Security

If you have a security issue to report, please contact us at [security@prisma.io](mailto:security@prisma.io?subject=[GitHub]%20Prisma%202%20Security%20Report%20Go)
