package types

import (
	"fmt"
	"testing"
)

func TestString_GoCase(t *testing.T) {
	tests := []struct {
		have String
		want string
	}{{
		have: "",
		want: "",
	}, {
		have: "anotherIDStuffSomethingID",
		want: "AnotherIDStuffSomethingID",
	}, {
		have: "anotherIdStuffSomethingId",
		want: "AnotherIDStuffSomethingID",
	}, {
		have: "APISession",
		want: "APISession",
	}}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.have, tt.want), func(t *testing.T) {
			if got := tt.have.GoCase(); got != tt.want {
				t.Errorf("GoCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
