package composite

import (
	"context"
	"fmt"
	"github.com/steebchen/prisma-client-go/runtime/extension"
	"github.com/steebchen/prisma-client-go/test"
	"log"
	"testing"
	"time"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestExtend(t *testing.T) {
	//t.Skip("not implemented yet")
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "disallow deletes",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			extended := client.Prisma.Extend(
				ext.OnAll(func(ctx context.Context, model Model, operation Operation, args Args, run Run) (*Result, error) {
					log.Printf("on all: model: %s, operation: %s", model, operation)

					switch operation {
					case extension.DeleteOne, extension.DeleteMany:
						return nil, fmt.Errorf("deletes are not allowed")
					}

					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "disallow mutations from non-logged in user",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			GetViewer := func(ctx context.Context) (*struct{}, bool) {
				return nil, true
			}
			extended := client.Prisma.Extend(
				ext.OnMutation(func(ctx context.Context, model Model, operation Operation, args Args, run Run) (*Result, error) {
					// adapt this to how you retrieve the current user from your HTTP context
					_, ok := GetViewer(ctx)
					if !ok {
						return nil, fmt.Errorf("not logged in")
					}
					log.Printf("on all: model: %s, operation: %s", model, operation)

					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "track timing",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			extended := client.Prisma.Extend(
				ext.OnAll(func(ctx context.Context, model Model, operation Operation, args Args, run Run) (*Result, error) {
					now := time.Now()
					defer func() {
						log.Printf("timing: %s %s %s", model, operation, time.Since(now))
					}()
					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "on operation",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			extended := client.Prisma.Extend(
				ext.OnOperation([]Operation{extension.FindUnique}, func(ctx context.Context, model Model, args Args, run Run) (*Result, error) {
					log.Printf("on operation: model: %s", model)
					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "on model",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			extended := client.Prisma.Extend(
				ext.OnModel([]Model{UserModelName}, func(ctx context.Context, operation Operation, args Args, run Run) (*Result, error) {
					log.Printf("on operation: operation: %s", operation)
					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "on model and operation",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			extended := client.Prisma.Extend(
				ext.On([]Operation{extension.FindUnique}, []Model{UserModelName}, func(ctx context.Context, model Model, operation Operation, args Args, run Run) (*Result, error) {
					return run()
				}),
			)
			log.Printf("extended: %+v", extended)
		},
	}, {
		name:   "individual extensions",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			t.Skipf("not implemented yet")
			//extended := client.Prisma.Extend(
			// TODO this needs a lot more work as all the params have to be adapted to the respective method
			//ext.OnUserDeleteMany(func(ctx context.Context, query UserWhereParam, run Run) (*Result, error) {
			//	return run()
			//}),
			//)
			//log.Printf("extended: %+v", extended)
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, test.Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
