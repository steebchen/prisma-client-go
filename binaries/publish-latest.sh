#!/bin/sh

set -eux

v=$(curl -s https://api.github.com/repos/prisma/prisma/releases/latest | jq -r .tag_name)

sh publish.sh "$v"
