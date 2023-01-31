package raw

import (
	"encoding/json"
	"testing"

	"github.com/prisma/prisma-client-go/test/helpers/massert"
)

func TestBytes_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Bytes
		args     args
		wantErr  bool
	}{{
		name:     "object",
		expected: []byte(`{"some":5}`),
		args: args{
			b: []byte(`{"prisma__type":"bytes","prisma__value":{"some":5}}`),
		},
	}, {
		name:     "string",
		expected: []byte(`"asdf"`),
		args: args{
			b: []byte(`{"prisma__type":"bytes","prisma__value":"asdf"}`),
		},
	}, {
		name:     "number",
		expected: []byte("5"),
		args: args{
			b: []byte(`{"prisma__type":"bytes","prisma__value":5}`),
		},
	}, {
		name:    "error on wrong type",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"string","prisma__value":5}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Bytes
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			massert.Equal(t, tt.expected, v)
		})
	}
}
