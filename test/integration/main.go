package main

import (
	"context"
	"fmt"

	"integration/db"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	client := db.NewClient()
	check(client.Prisma.Connect())
	defer func() {
		check(client.Prisma.Disconnect())
	}()

	ctx := context.Background()

	result, err := client.User.FindMany().Delete().Exec(ctx)
	check(err)
	fmt.Printf("remove %d items\n", result.Count)

	_, err = client.User.CreateOne(
		db.User.Email.Set("new@email.com"),
		db.User.Name.Set("John"),
	).Exec(ctx)
	check(err)

	user, err := client.User.FindUnique(
		db.User.Email.Equals("new@email.com"),
	).Exec(ctx)
	check(err)

	fmt.Printf("user: %+v\n", user)

	name, ok := user.Name()
	fmt.Printf("nullable name: %s %t\n", name, ok)
}
