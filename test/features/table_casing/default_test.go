package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestTableCasing(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		mockDB := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDB)

		expectedUser := UserModel{
			RawUser: RawUser{
				ID: "user",
			},
		}

		actualCreatedUser, err := client.User.CreateOne(
			User.ID.Set("user"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedUser, actualCreatedUser)

		actualFoundUser, err := client.User.FindOne(
			User.ID.Equals("user"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedUser, actualFoundUser)

		expectedEventLower := EventLowerModel{
			RawEventLower: RawEventLower{
				ID: "event",
			},
		}

		actualCreatedEventLower, err := client.EventLower.CreateOne(
			EventLower.ID.Set("event"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedEventLower, actualCreatedEventLower)

		actualFoundEventLower, err := client.EventLower.FindOne(
			EventLower.ID.Equals("event"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedEventLower, actualFoundEventLower)

		expectedParticipant := ParticipantUpperModel{
			RawParticipantUpper: RawParticipantUpper{
				ID:      "participant",
				EventID: "event",
			},
		}

		actualCreatedParticipant, err := client.ParticipantUpper.CreateOne(
			ParticipantUpper.Event.Link(
				EventLower.ID.Equals("event"),
			),
			ParticipantUpper.ID.Set("participant"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedParticipant, actualCreatedParticipant)

		actualFoundParticipant, err := client.ParticipantUpper.FindOne(
			ParticipantUpper.ID.Equals("participant"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedParticipant, actualFoundParticipant)
	})
}
