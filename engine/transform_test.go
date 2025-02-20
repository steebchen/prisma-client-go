package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		name: "transform",
		args: args{
			data: []byte(`{"columns":["id","email","username","str","strOpt","date","dateOpt","int","intOpt","float","floatOpt","bool","boolOpt"],"types":["string","string","string","string","string","datetime","datetime","int","int","double","double","int","int"],"rows":[["id1","email1","a","str","strOpt","2020-01-01T00:00:00+00:00","2020-01-01T00:00:00+00:00",5,5,5.5,5.5,1,0],["id2","email2","b","str","strOpt","2020-01-01T00:00:00+00:00","2020-01-01T00:00:00+00:00",5,5,5.5,5.5,1,0]]}`),
		},
		want: []byte(`[{"bool":1,"boolOpt":0,"date":"2020-01-01T00:00:00+00:00","dateOpt":"2020-01-01T00:00:00+00:00","email":"email1","float":5.5,"floatOpt":5.5,"id":"id1","int":5,"intOpt":5,"str":"str","strOpt":"strOpt","username":"a"},{"bool":1,"boolOpt":0,"date":"2020-01-01T00:00:00+00:00","dateOpt":"2020-01-01T00:00:00+00:00","email":"email2","float":5.5,"floatOpt":5.5,"id":"id2","int":5,"intOpt":5,"str":"str","strOpt":"strOpt","username":"b"}]`),
	},
		{
			name: "transform mongo raw response",
			args: args{
				data: []byte(`[{"_id":{"$oid":"67347ee4a18fa09750c1085a"},"createdAt":{"$date":"2024-11-13T10:26:44.246Z"},"firstName":"Trua Nguyen","lastName":"Van","email":"truanv@gmail"},{"_id":{"$oid":"67348094597e341917026845"},"email":"truanv@gmail","firstName":"Trua Nguyen","lastName":"Van"},{"_id":{"$oid":"673480d6597e341917026dea"},"email":"truanv@gmail","firstName":"Trua Nguyen","lastName":"Van"},{"_id":{"$oid":"67348265597e34191702904f"},"firstName":"Trua Nguyen ","lastName":"Van","email":"truanv@gmail"},{"_id":{"$oid":"6734827b597e34191702923d"},"email":"truanv@gmail","firstName":"Trua Nguyen ","lastName":"Van"}]`),
			},
			want: []byte(`[{"_id":"67347ee4a18fa09750c1085a","createdAt":"2024-11-13T10:26:44.246Z","email":"truanv@gmail","firstName":"Trua Nguyen","id":"67347ee4a18fa09750c1085a","lastName":"Van"},{"_id":"67348094597e341917026845","email":"truanv@gmail","firstName":"Trua Nguyen","id":"67348094597e341917026845","lastName":"Van"},{"_id":"673480d6597e341917026dea","email":"truanv@gmail","firstName":"Trua Nguyen","id":"673480d6597e341917026dea","lastName":"Van"},{"_id":"67348265597e34191702904f","email":"truanv@gmail","firstName":"Trua Nguyen ","id":"67348265597e34191702904f","lastName":"Van"},{"_id":"6734827b597e34191702923d","email":"truanv@gmail","firstName":"Trua Nguyen ","id":"6734827b597e34191702923d","lastName":"Van"}]`),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformResponse(tt.args.data)
			if err != nil {
				t.Fatalf("transformResponse() error = %v", err)
			}
			assert.Equal(t, string(tt.want), string(got))
		})
	}
}
