#!/bin/sh

set -eux

go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name init
go run github.com/prisma/prisma-client-go migrate up --experimental
