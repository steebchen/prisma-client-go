#!/bin/sh

set -eux

url="https://api.github.com/repos/prisma/prisma/releases/latest"

i=1
while [ $i -le 20 ]; do
  v=$(curl -s "$url" -H "Authorization: Bearer $GH_TOKEN" | jq -r .tag_name)
  if [ "$v" != "null" ]; then
    break
  fi
  sleep $i
  i=$((i+1))
done

if [ "$v" = "null" ]; then
  echo "Could not find latest version"
  echo "full response:"
  curl -s "$url" -H "Authorization: Bearer $GH_TOKEN"
  exit 1
fi

sh publish.sh "$v"
