# Prisma CLI

## Prisma CLI with Go

The go client works slightly different from the normal Prisma tooling. When you're using the go client, whenever you see Prisma CLI commands such as `prisma ...`, you should always write `go run github.com/steebchen/prisma-client-go ...` instead.

## Shortcut

If you just work with the Go client and don't have (or want) the NodeJS Prisma CLI installed, you can set up an alias so that you can write `prisma` commands as usual, but it'll invoke the real locally bundled Prisma CLI. To do that, edit your `~/.bashrc` or `~/.zshrc` and add:

```
alias prisma="go run github.com/steebchen/prisma-client-go"
```

Now `prisma generate` and any other command will work, and it'll just run 1`go run github.com/steebchen/prisma-client-go generate` under the hood.
