package raw

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parse(v string) time.Time {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		panic(err)
	}
	return t
}

func TestTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Time
		args     args
		wantErr  bool
	}{{
		name:     "date value",
		expected: Time{Time: parse("2023-05-05T05:05:05Z")},
		args: args{
			b: []byte(`{"prisma__type":"date","prisma__value":"2023-05-05T05:05:05Z"}`),
		},
	}, {
		name:    "date error on wrong data",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"date","prisma__value":5}`),
		},
	}, {
		name:     "datetime value",
		expected: Time{Time: parse("2023-05-05T05:05:05Z")},
		args: args{
			b: []byte(`{"prisma__type":"datetime","prisma__value":"2023-05-05T05:05:05Z"}`),
		},
	}, {
		name:    "datetime error on wrong data",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"datetime","prisma__value":5}`),
		},
	}, {
		name:    "error on wrong type",
		wantErr: true,
		args: args{
			b: []byte(`{"prisma__type":"string","prisma__value":"2023-01-01T00:00:00Z"}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Time
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
