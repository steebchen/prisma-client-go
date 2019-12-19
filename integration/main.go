package main

import (
	"context"
	"fmt"

	"github.com/prisma/photongo/integration/photon"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	client := photon.NewClient()
	err := client.Connect()
	check(err)
	defer func() {
		err := client.Disconnect()
		check(err)
	}()

	ctx := context.Background()

	count, err := client.User.FindMany().Delete().Exec(ctx)
	check(err)
	fmt.Printf("remove %d items\n", count)

	_, err = client.User.CreateOne(
		photon.User.Email.Set("new@email.com"),
		photon.User.Name.Set("John"),
	).Exec(ctx)
	check(err)

	user, err := client.User.FindOne(
		photon.User.Email.Equals("new@email.com"),
	).Exec(ctx)
	check(err)

	fmt.Printf("user: %+v\n", user)

	name, ok := user.Name()
	fmt.Printf("nullable name: %s %t\n", name, ok)
}
