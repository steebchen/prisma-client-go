# Contributing

## PRs Welcome

We welcome contributions to this project. For small things, feel free to directly create a PR. For major changes, please open an issue first to discuss the changes.

## Writing Commit Messages

We use [conventional commits](https://www.conventionalcommits.org) (also known as semantic commits) to ensure consistent and descriptive commit messages.

## Tests

### Running tests

```shell
# setup deps & generate code – requires docker to be installed
# this starts a docker compose stack with all required databases
go generate -tags setup ./...
# if you already ran setup, just run the following
go generate ./...
go test ./... -v

# to teardown docker containers:
go generate -tags teardown ./...
```

### How integration tests work

Most test live in the `test/` directory and are integration tests of the generated client. That means there's a Prisma
schema and before running the test, the client needs to be generated first. There may be table-driven tests which, on
each individual test run, creates a new isolated database, runs migrations, then run the tests, and finally cleans up
the database afterwards.

You can also run individual code generation tests via your editor, however keep in mind you need to run
`go generate ./...` before in the directory of the tests you want to run.

### E2E tests

End-to-end tests require third party credentials and may also be flaky from time to time. This is why they are not run locally by default and optional in CI.

To run them locally, you need to set up all required credentials (check the [env vars used for CI](https://github.com/steebchen/prisma-client-go/blob/a8a05c34aadd035303ea4651fcf6187cc4d039a0/.github/workflows/e2e-test.yml#L43), and then run:

```sh
cd test/e2e/
go generate -tags e2e ./...
go test ./... -run '^TestE2E.*$' -tags e2e -v
```
