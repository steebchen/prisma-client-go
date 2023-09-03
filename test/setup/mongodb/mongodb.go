package mongodb

import (
	"fmt"
	"testing"

	"github.com/steebchen/prisma-client-go/test/setup"
)

var MongoDB = &mongoDB{}

type mongoDB struct{}

func (*mongoDB) Name() string {
	return "mongodb"
}

func (*mongoDB) ConnectionString(mockDBName string) string {
	return fmt.Sprintf("mongodb://prisma:pw@localhost:27016/%s?authSource=admin&retryWrites=true", mockDBName)
}

func (*mongoDB) SetupDatabase(t *testing.T) string {
	return setup.RandomString()
}

func (*mongoDB) TeardownDatabase(t *testing.T, mockDB string) {}
