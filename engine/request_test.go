package engine

import (
	"reflect"
	"testing"
)

func Test_transformResponse(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{{
		name: "replace nulls",
		args: args{
			data: []byte(`[{"prisma__type":"string","prisma__value":"asdf"},{"prisma__type":"null","prisma__value":null}]`),
		},
		want: []byte(`[{"prisma__type":"string","prisma__value":"asdf"},null]`),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformResponse(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
