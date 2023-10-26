package raw

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client := NewClient(
		WithDatasourceURL("file:manual-override.sqlite"),
	)

	mockDB := test.Start(t, test.SQLite, client.Engine, []string{})
	defer test.End(t, test.SQLite, client.Engine, mockDB)

	user, err := client.User.FindUnique(
		User.ID.Equals("persisted-1"),
	).Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &UserModel{
		InnerUser: InnerUser{
			ID:    "persisted-1",
			Email: "persisted@example.com",
		},
		RelationsUser: RelationsUser{},
	}, user)

	// this code was used to create the entry and the test DB
	// user, err := client.User.CreateOne(
	// 	User.Email.Set("persisted@example.com"),
	// 	User.ID.Set("persisted-1"),
	// ).Exec(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
