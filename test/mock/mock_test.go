package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypedMock(t *testing.T) {
	do := func(ctx context.Context, client *PrismaClient) (UserModel, error) {
		user, err := client.User.FindOne(User.ID.Equals("foo")).Exec(ctx)
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

	client, mock, ensure := NewMock()
	defer ensure(t)
	mock.User.Expect(
		client.User.FindOne(User.ID.Equals("foo")),
	).Returns(expected)

	actual, err := do(context.Background(), client)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expected, actual)
}

func TestMockError(t *testing.T) {
	do := func(ctx context.Context, client *PrismaClient) (UserModel, error) {
		user, err := client.User.FindOne(User.ID.Equals("foo")).Exec(ctx)
		if err != nil {
			return UserModel{}, err
		}

		return user, nil
	}

	expectedErr := ErrNotFound
	expected := UserModel{}

	client, mock, ensure := NewMock()
	defer ensure(t)
	mock.User.Expect(
		client.User.FindOne(User.ID.Equals("foo")),
	).Errors(ErrNotFound)

	actual, err := do(context.Background(), client)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expected, actual)
}
