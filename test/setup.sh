#!/bin/sh

set -eux

cd test/setup/setup/
go generate -tags setup ./...

docker ps | grep go-client
