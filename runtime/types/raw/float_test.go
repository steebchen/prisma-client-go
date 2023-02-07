package raw

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Float
		args     args
		wantErr  bool
	}{{
		name:     "zero",
		expected: 0,
		args: args{
			b: []byte(`{"prisma__type":"double","prisma__value":0}`),
		},
	}, {
		name:     "value",
		expected: -5.3823923828,
		args: args{
			b: []byte(`{"prisma__type":"double","prisma__value":-5.3823923828}`),
		},
	}, {
		name: "string value",
		args: args{
			b: []byte(`{"prisma__type":"double","prisma__value":"123.456"}`),
		},
		wantErr: true,
	}, {
		name:    "error on wrong type",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"string","prisma__value":5}`),
		},
	}, {
		name:    "error on wrong data",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"double","prisma__value":true}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Float
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
