package mock

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	now := time.Now()
	do := func(ctx context.Context, db *PrismaClient) (UserModel, error) {
		user, err := db.User.FindOne(User.ID.Equals("foo")).Exec(ctx)
		if err != nil {
			return UserModel{}, err
		}

		return user, nil
	}

	var expectedErr error = nil
	expected := UserModel{
		InternalUser: InternalUser{
			ID:   "123",
			Name: "foo",
		},
		RelationsUser: RelationsUser{},
	}

	db, mock, ensure := NewMock()
	defer ensure(t)
	mock.Expect(
		db.User.FindOne(User.ID.Equals("foo")),
	).Returns(expected)

	actual, err := do(context.Background(), db)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expected, actual)

	log.Printf("%s", time.Since(now))
}
