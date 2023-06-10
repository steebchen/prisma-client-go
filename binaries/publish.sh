#!/bin/sh

set -eux

S3_BUCKET="prisma-photongo"

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

npx pkg -t node16-linux-x64,node16-darwin-x64,node16-win-x64,node16-linux-arm64,node16-darwin-arm64,node16-win-arm64 node_modules/prisma

version=$(npx prisma version | grep '^\(prisma \)' | cut -d : -f 2 | cut -d " " -f 2)
mv prisma-macos-x64 "prisma-cli-$version-darwin-x64"
mv prisma-linux-x64 "prisma-cli-$version-linux-x64"
mv prisma-win-x64.exe "prisma-cli-$version-windows-x64.exe"
mv prisma-macos-arm64 "prisma-cli-$version-darwin-arm64"
mv prisma-linux-arm64 "prisma-cli-$version-linux-arm64"
mv prisma-win-arm64.exe "prisma-cli-$version-windows-arm64.exe"

gzip "prisma-cli-$version-darwin-x64"
gzip "prisma-cli-$version-linux-x64"
gzip "prisma-cli-$version-windows-x64.exe"
gzip "prisma-cli-$version-darwin-arm64"
gzip "prisma-cli-$version-linux-arm64"
gzip "prisma-cli-$version-windows-arm64.exe"

aws s3 cp "prisma-cli-$version-darwin-x64.gz" "s3://$S3_BUCKET" --acl public-read
aws s3 cp "prisma-cli-$version-linux-x64.gz" "s3://$S3_BUCKET" --acl public-read
aws s3 cp "prisma-cli-$version-windows-x64.exe.gz" "s3://$S3_BUCKET" --acl public-read
aws s3 cp "prisma-cli-$version-darwin-arm64.gz" "s3://$S3_BUCKET" --acl public-read
aws s3 cp "prisma-cli-$version-linux-arm64.gz" "s3://$S3_BUCKET" --acl public-read
aws s3 cp "prisma-cli-$version-windows-arm64.exe.gz" "s3://$S3_BUCKET" --acl public-read

cd ../..

# cleanup
rm -rf build
