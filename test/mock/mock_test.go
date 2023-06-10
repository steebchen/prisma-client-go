package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypedMockReturns(t *testing.T) {
	do := func(ctx context.Context, client *PrismaClient) (*UserModel, error) {
		user, err := client.User.FindUnique(User.ID.Equals("foo")).Exec(ctx)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	var expectedErr error
	expected := &UserModel{
		InnerUser: InnerUser{
			ID:   "123",
			Name: "foo",
		},
		RelationsUser: RelationsUser{},
	}

	client, mock, ensure := NewMock()
	defer ensure(t)
	mock.User.Expect(
		client.User.FindUnique(User.ID.Equals("foo")),
	).Returns(*expected)

	actual, err := do(context.Background(), client)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expected, actual)
}

func TestTypedMockReturnsMany(t *testing.T) {
	do := func(ctx context.Context, client *PrismaClient) ([]UserModel, error) {
		return client.User.FindMany(User.Name.Equals("foo")).Exec(ctx)
	}

	var expectedErr error
	expected := []UserModel{
		{
			InnerUser: InnerUser{
				ID:   "123",
				Name: "foo",
			},
			RelationsUser: RelationsUser{},
		},
	}

	client, mock, ensure := NewMock()
	defer ensure(t)
	mock.User.Expect(
		client.User.FindMany(User.Name.Equals("foo")),
	).ReturnsMany(expected)

	actual, err := do(context.Background(), client)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, expected, actual)
}

func TestMockError(t *testing.T) {
	do := func(ctx context.Context, client *PrismaClient) (*UserModel, error) {
		user, err := client.User.FindUnique(User.ID.Equals("foo")).Exec(ctx)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	expectedErr := ErrNotFound

	client, mock, ensure := NewMock()
	defer ensure(t)
	mock.User.Expect(
		client.User.FindUnique(User.ID.Equals("foo")),
	).Errors(ErrNotFound)

	actual, err := do(context.Background(), client)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, true, actual == nil)
}
