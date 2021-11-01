#!/bin/sh

export PRISMA_CLIENT_ENGINE_TYPE=dataproxy
go run github.com/prisma/prisma-client-go generate --schema schema.out.prisma
