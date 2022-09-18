package main

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/engine/migrate"
)

const schmea = `
datasource db {
        // could be postgresql or mysql
        provider = "sqlite"
        url      = "file:dev.db"
    }

    generator db {
        provider = "go run github.com/prisma/prisma-client-go"
        // set the output folder and package name
        // output           = "./your-folder"
        // package          = "yourpackagename"
    }

    model Post {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        updatedAt DateTime @updatedAt
        title     String
        published Boolean
        desc      String?
    }
`

func main() {
	// if err := run(); err != nil {
	// 	panic(err)
	// }
	engine := migrate.NewMigrationEngine()

	engine.Push("schema.prisma")
	engine.Push2("schema.prisma")
	engine.Push("schema.prisma")
	engine.Push2("schema.prisma")

	// testDmmf()
	// testSdl()
}

func testDmmf() {
	engine := engine.NewQueryEngine(schmea, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	dmmf, err := engine.IntrospectDMMF(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(dmmf.Datamodel)
}

func testSdl() {
	engine := engine.NewQueryEngine(schmea, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	ctx := context.TODO()
	sdl, err := engine.IntrospectSDL(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(sdl))

}

// func run() error {
// 	client := db.NewClient()
// 	if err := client.Prisma.Connect(); err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err := client.Prisma.Disconnect(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	ctx := context.Background()

// 	// create a post
// 	createdPost, err := client.Post.CreateOne(
// 		db.Post.Title.Set("Hi from Prisma!"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, _ := json.MarshalIndent(createdPost, "", "  ")
// 	fmt.Printf("created post: %s\n", result)

// 	// find a single post
// 	post, err := client.Post.FindUnique(
// 		db.Post.ID.Equals(createdPost.ID),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, _ = json.MarshalIndent(post, "", "  ")
// 	fmt.Printf("post: %s\n", result)

// 	// for optional/nullable values, you need to check the function and create two return values
// 	// `desc` is a string, and `ok` is a bool whether the record is null or not. If it's null,
// 	// `ok` is false, and `desc` will default to Go's default values; in this case an empty string (""). Otherwise,
// 	// `ok` is true and `desc` will be "my description".
// 	desc, ok := post.Desc()
// 	if !ok {
// 		return fmt.Errorf("post's description is null")
// 	}

// 	fmt.Printf("The posts's description is: %s\n", desc)

// 	createUserA := client.Post.CreateOne(
// 		db.Post.Title.Set("2"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("222."),
// 	).Tx()

// 	createUserB := client.Post.CreateOne(
// 		db.Post.Title.Set("3"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("222."),
// 	).Tx()

// 	tx := client.Prisma.Transaction(createUserA, createUserB)
// 	if err := tx.Exec(ctx); err != nil {
// 		panic(err)
// 	}

// 	return nil
// }
