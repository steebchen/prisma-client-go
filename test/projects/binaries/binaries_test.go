package binaries

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/binaries/platform"
)

func TestBinaries(t *testing.T) {
	t.Parallel()

	// just for logging purposes
	out, _ := exec.Command("ls").CombinedOutput()
	log.Printf("%s", string(out))

	// this test only verifies that specifying `binaryTargets` downloaded the separate file into the directory
	_, err := os.Stat("./prisma-query-engine-" + platform.BinaryPlatformName())
	assert.Equal(t, err, nil)

	_, err = os.Stat("./prisma-query-engine-rhel-openssl-1.1.x")
	assert.Equal(t, err, nil)
}
