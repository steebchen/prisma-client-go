package raw

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Int
		args     args
		wantErr  bool
	}{{
		name:     "zero",
		expected: 0,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":0}`),
		},
	}, {
		name:     "value",
		expected: -5,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":-5}`),
		},
	}, {
		name:    "string value",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"int","prisma__value":"4"}`),
		},
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
			b: []byte(`{"prisma__type":"int","prisma__value":true}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Int
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
