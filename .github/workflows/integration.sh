#!/bin/sh

set -eux

# create network
docker network create --driver nat integration-bridge

# launch postgres
docker run --network integration-bridge --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres

# run migrations
docker build . -f docker/migrate.dockerfile -t migrate
docker run --network integration-bridge migrate

# run actual tests
docker build . -f docker/integration.dockerfile -t integration
docker run --network integration-bridge integration

# cleanup
docker stop postgres
docker rm postgres
docker network remove integration-bridge
