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
		name: "replace nulls array",
		args: args{
			data: []byte(`[{"prisma__type":"string","prisma__value":"asdf"},{"prisma__type":"null","prisma__value":null}]`),
		},
		want: []byte(`["asdf",null]`),
	}, {
		name: "replace nulls object",
		args: args{
			data: []byte(`[{"id":{"prisma__type":"null","prisma__value":null}}]`),
		},
		want: []byte(`[{"id":null}]`),
	}, {
		name: "replace string",
		args: args{
			data: []byte(`[{"id":{"prisma__type":"string","prisma__value":"asdf"}}]`),
		},
		want: []byte(`[{"id":"asdf"}]`),
	}, {
		name: "replace string with other objects",
		args: args{
			data: []byte(`[{"id":{"prisma__type":"string","prisma__value":"asdf"}},{"some":{"other":"object"}},{"value":5}]`),
		},
		want: []byte(`[{"id":"asdf"},{"some":{"other":"object"}},{"value":5}]`),
	}, {
		name: "native string",
		args: args{
			data: []byte(`"asdf"`),
		},
		want: []byte(`"asdf"`),
	}, {
		name: "native number",
		args: args{
			data: []byte(`5`),
		},
		want: []byte(`5`),
	}, { // edge cases which are specifically handled
		name: "bytes",
		args: args{
			data: []byte(`{"item":{"prisma__type":"bytes","prisma__value":"eyJzb21lIjo1fQ=="}}`),
		},
		want: []byte(`{"item":"eyJzb21lIjo1fQ=="}`),
	}, {
		name: "bytes",
		args: args{
			data: []byte(`[{"prisma__type":"bytes","prisma__value":"eyJzb21lIjp7ImEiOiJiIn19"}]`),
		},
		want: []byte(`["eyJzb21lIjp7ImEiOiJiIn19"]`),
	}, {
		name: "bytes",
		args: args{
			data: []byte(`[{"prisma__type":"bytes","prisma__value":"MTIz"}]`),
		},
		want: []byte(`["MTIz"]`),
	}, {
		name: "big",
		args: args{
			data: []byte(`[{"id":{"prisma__type":"string","prisma__value":"id1"},"email":{"prisma__type":"string","prisma__value":"email1"},"username":{"prisma__type":"string","prisma__value":"a"},"str":{"prisma__type":"string","prisma__value":"str"},"strOpt":{"prisma__type":"string","prisma__value":"strOpt"},"strEmpty":{"prisma__type":"null","prisma__value":null},"time":{"prisma__type":"datetime","prisma__value":"2020-01-01T00:00:00+00:00"},"timeOpt":{"prisma__type":"datetime","prisma__value":"2020-01-01T00:00:00+00:00"},"timeEmpty":{"prisma__type":"null","prisma__value":null},"int":{"prisma__type":"int","prisma__value":5},"intOpt":{"prisma__type":"int","prisma__value":5},"intEmpty":{"prisma__type":"null","prisma__value":null},"float":{"prisma__type":"double","prisma__value":5.5},"floatOpt":{"prisma__type":"double","prisma__value":5.5},"floatEmpty":{"prisma__type":"null","prisma__value":null},"bool":{"prisma__type":"bool","prisma__value":true},"boolOpt":{"prisma__type":"bool","prisma__value":false},"boolEmpty":{"prisma__type":"null","prisma__value":null},"decimal":{"prisma__type":"decimal","prisma__value":"5.5"},"decimalOpt":{"prisma__type":"decimal","prisma__value":"5.5"},"decimalEmpty":{"prisma__type":"null","prisma__value":null},"json":{"prisma__type":"json","prisma__value":{"field":"value"}},"jsonOpt":{"prisma__type":"json","prisma__value":{"field":"value"}},"jsonEmpty":{"prisma__type":"null","prisma__value":null},"bytes":{"prisma__type":"bytes","prisma__value":"eyJmaWVsZCI6InZhbHVlIn0="},"bytesOpt":{"prisma__type":"bytes","prisma__value":"eyJmaWVsZCI6InZhbHVlIn0="},"bytesEmpty":{"prisma__type":"null","prisma__value":null}},{"id":{"prisma__type":"string","prisma__value":"id2"},"email":{"prisma__type":"string","prisma__value":"email2"},"username":{"prisma__type":"string","prisma__value":"b"},"str":{"prisma__type":"string","prisma__value":"str"},"strOpt":{"prisma__type":"string","prisma__value":"strOpt"},"strEmpty":{"prisma__type":"null","prisma__value":null},"time":{"prisma__type":"datetime","prisma__value":"2020-01-01T00:00:00+00:00"},"timeOpt":{"prisma__type":"datetime","prisma__value":"2020-01-01T00:00:00+00:00"},"timeEmpty":{"prisma__type":"null","prisma__value":null},"int":{"prisma__type":"int","prisma__value":5},"intOpt":{"prisma__type":"int","prisma__value":5},"intEmpty":{"prisma__type":"null","prisma__value":null},"float":{"prisma__type":"double","prisma__value":5.5},"floatOpt":{"prisma__type":"double","prisma__value":5.5},"floatEmpty":{"prisma__type":"null","prisma__value":null},"bool":{"prisma__type":"bool","prisma__value":true},"boolOpt":{"prisma__type":"bool","prisma__value":false},"boolEmpty":{"prisma__type":"null","prisma__value":null},"decimal":{"prisma__type":"decimal","prisma__value":"5.5"},"decimalOpt":{"prisma__type":"decimal","prisma__value":"5.5"},"decimalEmpty":{"prisma__type":"null","prisma__value":null},"json":{"prisma__type":"json","prisma__value":{"field":"value"}},"jsonOpt":{"prisma__type":"json","prisma__value":{"field":"value"}},"jsonEmpty":{"prisma__type":"null","prisma__value":null},"bytes":{"prisma__type":"bytes","prisma__value":"eyJmaWVsZCI6InZhbHVlIn0="},"bytesOpt":{"prisma__type":"bytes","prisma__value":"eyJmaWVsZCI6InZhbHVlIn0="},"bytesEmpty":{"prisma__type":"null","prisma__value":null}}]`),
		},
		want: []byte(`[{"bool":true,"boolEmpty":null,"boolOpt":false,"bytes":"eyJmaWVsZCI6InZhbHVlIn0=","bytesEmpty":null,"bytesOpt":"eyJmaWVsZCI6InZhbHVlIn0=","decimal":"5.5","decimalEmpty":null,"decimalOpt":"5.5","email":"email1","float":5.5,"floatEmpty":null,"floatOpt":5.5,"id":"id1","int":5,"intEmpty":null,"intOpt":5,"json":{"field":"value"},"jsonEmpty":null,"jsonOpt":{"field":"value"},"str":"str","strEmpty":null,"strOpt":"strOpt","time":"2020-01-01T00:00:00+00:00","timeEmpty":null,"timeOpt":"2020-01-01T00:00:00+00:00","username":"a"},{"bool":true,"boolEmpty":null,"boolOpt":false,"bytes":"eyJmaWVsZCI6InZhbHVlIn0=","bytesEmpty":null,"bytesOpt":"eyJmaWVsZCI6InZhbHVlIn0=","decimal":"5.5","decimalEmpty":null,"decimalOpt":"5.5","email":"email2","float":5.5,"floatEmpty":null,"floatOpt":5.5,"id":"id2","int":5,"intEmpty":null,"intOpt":5,"json":{"field":"value"},"jsonEmpty":null,"jsonOpt":{"field":"value"},"str":"str","strEmpty":null,"strOpt":"strOpt","time":"2020-01-01T00:00:00+00:00","timeEmpty":null,"timeOpt":"2020-01-01T00:00:00+00:00","username":"b"}]`),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformResponse(tt.args.data)
			if err != nil {
				t.Fatalf("transformResponse() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformResponse() = %s, want %s", got, tt.want)
			}
		})
	}
}
