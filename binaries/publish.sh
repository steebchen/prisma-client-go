#!/bin/sh

set -eux

v="$1"

mkdir -p build
cd build
npm init --yes
npm i "pkg" --dev
npm i "prisma@$v" --dev
npm i "@prisma/client@$v"
npx prisma version

mkdir -p node_modules/prisma/node_modules/@prisma/engines
cp -R node_modules/@prisma/engines/* node_modules/prisma/node_modules/@prisma/engines

npx pkg -t node12-linux,node12-darwin,node12-win node_modules/prisma

version=$(npx prisma version | grep '^\(prisma \)' | cut -d : -f 2 | cut -d " " -f 2)
mv prisma-macos "prisma-cli-$version-darwin"
mv prisma-linux "prisma-cli-$version-linux"
mv prisma-win.exe "prisma-cli-$version-windows.exe"

gzip "prisma-cli-$version-darwin"
gzip "prisma-cli-$version-linux"
gzip "prisma-cli-$version-windows.exe"

aws s3 cp "prisma-cli-$version-darwin.gz" s3://prisma-photongo --acl public-read
aws s3 cp "prisma-cli-$version-linux.gz" s3://prisma-photongo --acl public-read
aws s3 cp "prisma-cli-$version-windows.exe.gz" s3://prisma-photongo --acl public-read

cd ../..

# cleanup
rm -rf build
