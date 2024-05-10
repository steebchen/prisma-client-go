#!/bin/sh

set -eux

v=$(curl -s https://api.github.com/repos/prisma/prisma/releases/latest | jq -r .tag_name)
for i in {1..20}; do
  if [ "$v" != "null" ]; then
    break
  fi
  v=$(curl -s https://api.github.com/repos/prisma/prisma/releases/latest | jq -r .tag_name)
  sleep $i
done

sh publish.sh "$v"
