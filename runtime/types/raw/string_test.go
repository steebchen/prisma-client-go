package raw

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected String
		args     args
		wantErr  bool
	}{{
		name:     "value",
		expected: "asdf",
		args: args{
			b: []byte(`{"prisma__type":"string","prisma__value":"asdf"}`),
		},
	}, {
		name:    "error on wrong type",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":"asdf"}`),
		},
	}, {
		name:    "error on wrong data",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"string","prisma__value":5"}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v String
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
