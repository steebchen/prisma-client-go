#!/bin/sh

set -eux

# run migrations
#docker build . -f docker/migrate.dockerfile -t migrate
#docker run -v "$(pwd)/integration/tmp/dev.db:/app" migrate

# run actual tests
docker build . -f docker/integration.dockerfile -t integration
docker run integration
