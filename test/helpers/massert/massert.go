package massert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Equal compares two objects via its MarshalJSON representation in tests
// This is useful for types using structs like Time or Decimal, which are
// hard to compare when using traditional assert.Equal
//
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
