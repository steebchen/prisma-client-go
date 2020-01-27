#!/bin/sh

set -eux

go run github.com/prisma/prisma-client-go lift save --create-db --name init
go run github.com/prisma/prisma-client-go lift up
