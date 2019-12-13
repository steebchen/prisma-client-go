# Prisma binaries

## How to build Prisma CLI binaries

### Setup

[Install node-packer](https://github.com/pmq20/node-packer)

I chose this package because [zeit/pkg](https://github.com/zeit/pkg) and [nexe/nexe](https://github.com/nexe/nexe) resulted in some weird errors.

### Build the binary

```shell script
# install prisma
npm i -g prisma@alpha
# build the binary
./nodec /usr/local/lib/node_modules/prisma2/build/index.js --skip-npm-install -o ./prisma-cli-linux-$(prisma2 -v | cut -f1 -d"," | sed 's/.*@//g')
# now, manually upload binary to s3 bucket `prisma-binaries-photongo`
# then, adapt the PRISMA_VERSION variable in binaries/binaries.go
```
