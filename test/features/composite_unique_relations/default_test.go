package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestComposite(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		// language=GraphQL
		mockDB := test.Start(t, db, client.Engine, []string{`
			mutation {
				createOneUser(data: {
					id: "user",
				}) {
					id
				}
			}
		`, `
			mutation {
				createOneEvent(data: {
					id: "event",
				}) {
					id
				}
			}
		`})
		defer test.End(t, db, client.Engine, mockDB)

		expectedParticipant := ParticipantModel{
			InternalParticipant: InternalParticipant{
				ID:      "new-participant",
				UserID:  "user",
				EventID: "event",
			},
		}

		actualCreatedParticipant, err := client.Participant.CreateOne(
			Participant.User.Link(
				User.ID.Equals("user"),
			),
			Participant.Event.Link(
				Event.ID.Equals("event"),
			),
			Participant.ID.Set("new-participant"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedParticipant, actualCreatedParticipant)

		actualFoundParticipant, err := client.Participant.FindOne(
			Participant.UserIDEventID(
				Participant.UserID.Equals("user"),
				Participant.EventID.Equals("event"),
			),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedParticipant, actualFoundParticipant)
	})
}
