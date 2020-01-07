#!/bin/sh

set -eux

go run github.com/prisma/photongo lift save --create-db --name init
go run github.com/prisma/photongo lift up
