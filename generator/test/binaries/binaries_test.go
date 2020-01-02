package binaries

//go:generate go run github.com/prisma/photongo generate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/photongo/binaries/platform"
)

func TestBinaries(t *testing.T) {
	t.Parallel()

	// this test only verifies that specifying `binaryTargets` downloaded the separate file into the directory
	_, err := os.Stat("./prisma-query-engine-" + platform.BinaryNameWithSSL())
	assert.Equal(t, err, nil)

	_, err = os.Stat("./prisma-query-engine-rhel-openssl-1.1.x")
	assert.Equal(t, err, nil)
}
