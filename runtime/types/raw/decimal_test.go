package raw

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimal_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		expected Decimal
		args     args
		wantErr  bool
	}{{
		name:     "zero",
		expected: Decimal{Decimal: decimal.NewFromFloat(0)},
		args: args{
			b: []byte(`{"prisma__type":"decimal","prisma__value":0}`),
		},
	}, {
		name:     "value",
		expected: Decimal{Decimal: decimal.NewFromFloat(-5.3823923828)},
		args: args{
			b: []byte(`{"prisma__type":"decimal","prisma__value":-5.3823923828}`),
		},
	}, {
		name:     "string value",
		expected: Decimal{Decimal: decimal.NewFromFloat(123.456)},
		args: args{
			b: []byte(`{"prisma__type":"decimal","prisma__value":"123.456"}`),
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
			b: []byte(`{"prisma__type":"decimal","prisma__value":true}`),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v Decimal
			if err := json.Unmarshal(tt.args.b, &v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, v)
		})
	}
}
