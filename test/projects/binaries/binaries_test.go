package binaries

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaries(t *testing.T) {
	t.Parallel()

	// check if specified engine in schema.prisma exists
	_, err := os.Stat("./query-engine-debian-openssl-1.1.x_gen.go")
	assert.Equal(t, err, nil)
}
