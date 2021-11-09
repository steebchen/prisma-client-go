# Prisma Data Proxy

Connection limit issues got you down? By using Prisma's Data Proxy, you can pool your connections to avoid overloading your database. You can also make use of a web data browser and builtin query console at [our Cloud platform](https://cloud.prisma.io).

![](https://user-images.githubusercontent.com/5013932/141067103-55ac8326-ddc1-4ad2-baac-dc9b916c9056.jpg)

In Prisma Cloud, navigate to the data proxy, enable it, and generate connection string which we will need later. This contains an API key so handle it like database credentials.

![](https://user-images.githubusercontent.com/5013932/141067009-56a088d6-8e4f-4508-81d8-eddc1df3943b.jpg)

Then, enable the data proxy preview feature in your prisma.schema:

```diff
generator client {
  provider        = "go run github.com/prisma/prisma-client-go"
+ previewFeatures = ["dataProxy"]
}
```

You can generate a "Data Proxy"-enabled Go Client in the Prisma schema or individually using an env var.
If you want to use the data proxy everywhere, you might want to set it in the schema, whereas when you just
want to use the data proxy in production, you might want to set an env var in your deployment script.

Setting it in the schema:

```diff
generator client {
  provider        = "go run github.com/prisma/prisma-client-go"
  previewFeatures = ["dataProxy"]
+ engineType      = "dataproxy"
}
```

Setting it just once when generating the client:

```
PRISMA_CLIENT_ENGINE_TYPE=dataproxy go run github.com/prisma/prisma-client-go generate
```

We recommend using an env var in your data source config, so that you can supply the previously generated connection string at runtime:

```prisma
datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}
```

Then you can use the Go client as you would normally. There won't be any generated query engine files as opposed to the normal mode.

```shell
export DATABASE_URL='prisma://.......' # set your previously generated data proxy connection string here
go run . # run your app
```

For more information and caveats, read the full [Prisma Data Proxy docs](https://www.prisma.io/docs/concepts/components/prisma-data-platform#prisma-data-proxy).
