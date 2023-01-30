package massert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//goland:noinspection GoUnusedExportedFunction
func Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	e, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	a, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(e), string(a), msgAndArgs...)
}
