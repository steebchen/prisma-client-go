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
    <a href="./docs/reference.md">API Reference</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://www.prisma.io/blog">Blog</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://slack.prisma.io/">Slack</a>
    <span>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
    <a href="https://twitter.com/prisma">Twitter</a>
</div>

<hr>

Prisma Client Go is an **auto-generated query builder** that enables **type-safe** database access and **reduces boilerplate**. You can use it as an alternative to traditional ORMs such as GOORM, sqlboiler and most database-specific tools.

It is part of the [Prisma](https://www.prisma.io/) ecosystem. Prisma provides database tools for data access, declarative data modeling, schema migrations and visual data management.

*NOTE*: Prisma Client Go is currently considered alpha software. There will be documented breaking changes with new [releases](https://github.com/prisma/prisma-client-go/releases).

## Getting started

To get started, [read our quickstart tutorial](./docs/quickstart.md) to add Prisma to your project in just a few minutes.

You also might want to read [deployment tips](./docs/deploy.md) and the [full API reference](./docs/reference.md).

## Notes

The go client works slightly different than the normal Prisma tooling. When you're using the go client, whenever you see Prisma CLI commands such as `prisma ...`, you should always write `go run github.com/prisma/prisma-client-go ...` instead.
