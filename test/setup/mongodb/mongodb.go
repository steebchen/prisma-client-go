package mongodb

import (
	"fmt"
	"log"
	"testing"

	"github.com/prisma/prisma-client-go/test/cmd"
	"github.com/prisma/prisma-client-go/test/setup"
)

const containerName = "go-client-mongodb"
const image = "prismagraphql/mongo-single-replica:4.4.3-bionic"

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

func (db *mongoDB) Setup() {
	if err := cmd.Run("docker", "run", "--name", containerName, "-p", "27016:27016", "-e", "MONGO_INITDB_ROOT_USERNAME=prisma", "-e", "MONGO_INITDB_ROOT_PASSWORD=pw", "-e", "MONGO_PORT=27016", "-d", image); err != nil {
		panic(err)
	}
}

func (db *mongoDB) Teardown() {
	if err := cmd.Run("docker", "stop", containerName); err != nil {
		log.Println(err)
	}

	if err := cmd.Run("docker", "rm", containerName, "--force"); err != nil {
		log.Println(err)
	}
}
