#!/bin/sh

set -eux

url="https://api.github.com/repos/prisma/prisma/releases/latest"

v=$(curl -s "$url" -H "Authorization: Bearer $GH_TOKEN" | jq -r .tag_name)
for i in {1..20}; do
  if [ "$v" != "null" ]; then
    break
  fi
  v=$(curl -s "$url" -H "Authorization: Bearer $GH_TOKEN" | jq -r .tag_name)
  sleep $i
done

if [ "$v" = "null" ]; then
  echo "Could not find latest version"
  echo "full response:"
  curl -s "$url" -H "Authorization: Bearer $GH_TOKEN"
  exit 1
fi

sh publish.sh "$v"
