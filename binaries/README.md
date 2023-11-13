# Prisma binaries

There are two types of binaries needed by the Go client, one being the Prisma Engine binaries and the other being the
Prisma CLI binaries.

Prisma Engine binaries are fully managed, maintained and automatically updated by the Prisma team, as they are also
needed for our NodeJS client.

Prisma CLI binaries are not officially managed and were just by the maintainers of the Go client. This is why there is a
some documentation here and a script on how to build, upload and bump the Prisma CLI binaries.

**NOTE: This is just for documentation purposes. The Prisma CLI
is [automatically published](https://github.com/steebchen/prisma-client-go/blob/main/.github/workflows/publish-cli.yml).
**

--------

## How to build Prisma CLI binaries

Prisma CLI binaries are automatically published to S3 by a GitHub action. You can follow the instructions below to build
these binaries yourself.

### Prerequisites

Requires NodeJS.

Install the [AWS CLI](https://aws.amazon.com/cli/) and authenticate.

### Build the binary and upload to S3

#### Publish the latest Prisma version

```shell script
sh publish-latest.sh
```

#### Publish a specific Prisma version

```shell script
sh publish.sh <version>
# e.g.
sh publish.sh 3.0.0
```

You can check the available versions on the [Prisma releases page](https://github.com/prisma/prisma/releases).

**NOTE**:

#### Prisma employees

Any Prisma employee can authenticate with the Prisma Go client account. If you are a community member and would like to
bump the binaries, please ask us to do so in the #prisma-client-go channel in our public Slack.

#### Community members

If you want to set up Prisma CLI binaries yourself, authenticate with your own AWS account and adapt the bucket name
in `publish.sh`.
When using the client, you will need to override the URL with env vars whenever you run the Go client, specifically
`PRISMA_CLI_URL` and `PRISMA_ENGINE_URL`. You can see the shape of these values
in [binaries/version.go#L3-L8](https://github.com/steebchen/prisma-client-go/blob/main/binaries/version.go#L3-L8).

This will also print the query engine version which you will need in the next step.

### Bump the binaries in the Go client

Go to `binaries/version.go` and adapt
the [`PrismaVersion`](https://github.com/steebchen/prisma-client-go/blob/main/binaries/version.go#L4)
and [`EngineVersion`](https://github.com/steebchen/prisma-client-go/blob/main/binaries/version.go#L8)
to the new version values.
Push to a new branch, create a PR, and merge if tests are green (
e.g. [#709](https://github.com/steebchen/prisma-client-go/pull/709)).

When internal breaking changes happen, adaptions may be needed.
