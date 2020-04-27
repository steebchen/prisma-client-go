package dmmf

import (
	"testing"

	"github.com/prisma/prisma-client-go/generator/types"
)

func TestDocument_RelationName(t *testing.T) {
	type fields struct {
		Datamodel Datamodel
		Schema    Schema
		Mappings  []Mapping
	}
	type args struct {
		from types.StringLike
		to   types.StringLike
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.String
	}{{
		name: "user posts",
		fields: fields{
			Schema: Schema{
				OutputTypes: []OutputType{{
					Name: "Post",
					Fields: []SchemaField{{
						Name: "author",
						OutputType: SchemaOutputType{
							Type: "User",
						},
					}},
					IsEmbedded: false,
				}, {
					Name: "User",
					Fields: []SchemaField{{
						Name: "posts",
						OutputType: SchemaOutputType{
							Type: "Post",
						},
					}},
					IsEmbedded: false,
				}},
			},
		},
		args: args{
			to:   types.String("Post"),
			from: types.Type("User"),
		},
		want: "posts",
	}, {
		name: "post author",
		fields: fields{
			Schema: Schema{
				OutputTypes: []OutputType{{
					Name: "Post",
					Fields: []SchemaField{{
						Name: "author",
						OutputType: SchemaOutputType{
							Type: "User",
						},
					}},
					IsEmbedded: false,
				}, {
					Name: "User",
					Fields: []SchemaField{{
						Name: "posts",
						OutputType: SchemaOutputType{
							Type: "Post",
						},
					}},
					IsEmbedded: false,
				}},
			},
		},
		args: args{
			from: types.String("Post"),
			to:   types.Type("User"),
		},
		want: "author",
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{
				Datamodel: tt.fields.Datamodel,
				Schema:    tt.fields.Schema,
				Mappings:  tt.fields.Mappings,
			}
			if got := d.ReverseRelationName(tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("ReverseRelationName() = %v, want %v", got, tt.want)
			}
		})
	}
}
