package binaries

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/binaries/platform"
)

func TestBinaries(t *testing.T) {
	t.Parallel()

	// this test only verifies that specifying `binaryTargets` downloaded the separate file into the directory
	_, err := os.Stat("./query-engine-" + platform.BinaryPlatformName() + "_gen.go")
	assert.Equal(t, err, nil)

	_, err = os.Stat("./query-engine-debian-openssl-1.1.x_gen.go")
	assert.Equal(t, err, nil)
}
