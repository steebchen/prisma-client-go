# Prisma binaries

## How to build Prisma CLI binaries

### Setup

Install [zeit/pkg](https://github.com/zeit/pkg):

```shell script
npm i -g pkg
```

### Build the binary and upload to S3

```shell script
mkdir -p build
cd build
npm init --yes
npm i @prisma/cli@latest

mkdir -p binaries

pkg node_modules/@prisma/cli --out-path binaries/

cd binaries

version=$(npx prisma version | grep '^\(prisma2\|@prisma/cli\)' | cut -d : -f 2 | cut -d " " -f 2)
mv cli-macos "prisma-cli-$version-darwin"
mv cli-linux "prisma-cli-$version-linux"
mv cli-win.exe "prisma-cli-$version-windows.exe"

gzip "prisma-cli-$version-darwin"
gzip "prisma-cli-$version-linux"
gzip "prisma-cli-$version-windows.exe"

aws s3 cp "prisma-cli-$version-darwin.gz" s3://prisma-photongo --acl public-read
aws s3 cp "prisma-cli-$version-linux.gz" s3://prisma-photongo --acl public-read
aws s3 cp "prisma-cli-$version-windows.exe.gz" s3://prisma-photongo --acl public-read
```
