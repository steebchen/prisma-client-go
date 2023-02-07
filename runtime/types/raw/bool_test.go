package raw

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Boolean
		args     args
		wantErr  bool
	}{{
		name:     "true",
		expected: true,
		args: args{
			b: []byte(`{"prisma__type":"bool","prisma__value":true}`),
		},
	}, {
		name:     "false",
		expected: false,
		args: args{
			b: []byte(`{"prisma__type":"bool","prisma__value":false}`),
		},
	}, {
		name:     "int 1",
		expected: true,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":1}`),
		},
	}, {
		name:     "int 0",
		expected: false,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":0}`),
		},
	}, {
		name:    "error on wrong type",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"double","prisma__value":0}`),
		},
	}, {
		name:    "error on wrong data",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":"0"}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Boolean
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
