package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestComposite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create unchecked scalar",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "user",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneEvent(data: {
					id: "event",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			t.Skip()
			expectedParticipant := &ParticipantModel{
				InnerParticipant: InnerParticipant{
					ID:      "new-participant",
					UserID:  "user",
					EventID: "event",
				},
			}

			actualCreatedParticipant, err := client.Participant.CreateOne(
				// TODO unchecked scalars don't compile
				// Participant.UserID.Set("user"),
				nil,
				// Participant.EventID.Set("event"),
				nil,
				Participant.ID.Set("new-participant"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedParticipant, actualCreatedParticipant)
		},
	}, {
		name: "create normal connect",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "user",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneEvent(data: {
					id: "event",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expectedParticipant := &ParticipantModel{
				InnerParticipant: InnerParticipant{
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
		},
	}, {
		name: "find unique by named key",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "user",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneEvent(data: {
					id: "event",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneParticipant(data: {
					id: "participant",
					userId: "user",
					eventId: "event",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expectedParticipant := &ParticipantModel{
				InnerParticipant: InnerParticipant{
					ID:      "participant",
					UserID:  "user",
					EventID: "event",
				},
			}

			actualFoundParticipant, err := client.Participant.FindUnique(
				Participant.MyCustomKey(
					Participant.UserID.Equals("user"),
					Participant.EventID.Equals("event"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedParticipant, actualFoundParticipant)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.SQLite, test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
