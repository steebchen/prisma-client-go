# Composite

## Composite queries

The examples use the following prisma schema:

```prisma
model Organization {
  platformId   String
  platformKind String

  name String

  repositories Repository[]

  @@unique(name: "organization_name", [platformId, name])

  // Compose the primary key with the platformKind and platformId fields
  @@id(name: "organization_id", [platformKind, platformId])
}

model Repository {
  platformId   String
  platformKind String

  orgId String?
  org   Organization? @relation(fields: [platformKind, orgId], references: [platformKind, platformId])

  // Compose the primary key with the platformKind and platformId fields
  @@id(name: "repository_id", [platformKind, platformId])
}
```

## Composite keys

You can query for composite keys using the composite key id, whether it is a unique or primary key.

```go
// query by primary key
org, err := client.Organization.FindUnique(
  Organization.OrganizationID( // @@id(name: "organization_id", ...) maps to .OrganizationID
    Organization.PlatformKind.Equals("private"),
    Organization.PlatformID.Equals("123"),
  ),
).Exec(ctx)
```

```go
// query by unique key
org, err := client.Organization.FindUnique(
  Organization.OrganizationName( // @@id(name: "organization_name", ...) maps to .OrganizationName
    Organization.PlatformID.Equals("test"),
    Organization.Name.Equals("test"),
  ),
).Exec(ctx)
```

## Create with composite primary keys

To create records with a composite primary key, just specify the fields in the correct order. You don't have to
explicitly specify the composite key name.

```go
org, err = client.Organization.CreateOne(
  Organization.PlatformID.Set("123"),
  Organization.PlatformKind.Set("private"),
  Organization.Name.Set("test"),
).Exec(ctx)
```
