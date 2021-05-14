package composite

import (
	"context"
	"testing"
	"time"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestComposite(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "self unchecked scalar",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			t.Skip()

			_, err := client.Document.CreateOne(
				Document.ID.Set("doc-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Document.CreateOne(
				Document.ID.Set("doc-2"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Event.CreateOne(
				Event.ID.Set("event-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.EventInstance.CreateOne(
				EventInstance.Event.Link(Event.ID.Equals("event-1")),
				EventInstance.Start.Set(time.Now()),
				EventInstance.End.Set(time.Now()),
				EventInstance.Summary.Link(Document.ID.Equals("doc-1")),
				EventInstance.ID.Set("event-instance-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.EventInstance.CreateOne(
				EventInstance.Event.Link(Event.ID.Equals("event-1")),
				EventInstance.Start.Set(time.Now()),
				EventInstance.End.Set(time.Now()),
				EventInstance.Summary.Link(Document.ID.Equals("doc-2")),
				EventInstance.ID.Set("event-instance-2"),
				EventInstance.PreviousEventInstanceID.Set("event-instance-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

		},
	}, {
		name:   "self link",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.Document.CreateOne(
				Document.ID.Set("doc-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Document.CreateOne(
				Document.ID.Set("doc-2"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Event.CreateOne(
				Event.ID.Set("event-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.EventInstance.CreateOne(
				EventInstance.Event.Link(Event.ID.Equals("event-1")),
				EventInstance.Start.Set(time.Now()),
				EventInstance.End.Set(time.Now()),
				EventInstance.Summary.Link(Document.ID.Equals("doc-1")),
				EventInstance.ID.Set("event-instance-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.EventInstance.CreateOne(
				EventInstance.Event.Link(Event.ID.Equals("event-1")),
				EventInstance.Start.Set(time.Now()),
				EventInstance.End.Set(time.Now()),
				EventInstance.Summary.Link(Document.ID.Equals("doc-2")),
				EventInstance.Previous.Link(
					EventInstance.ID.Equals("event-instance-1"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}
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
