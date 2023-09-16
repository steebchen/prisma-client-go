#!/bin/sh

set -eux

cd test/setup/teardown/
go generate -tags teardown ./...

docker ps | grep go-client || true
