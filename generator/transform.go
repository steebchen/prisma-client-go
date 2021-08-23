package generator

import (
	"encoding/json"
	"fmt"
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
	"github.com/prisma/prisma-client-go/generator/ast/transform"
	"github.com/prisma/prisma-client-go/generator/types"
	"log"
	"os"
)

// Transform builds the AST from the flat DMMF so it can be used properly in templates
func Transform(input *Root) {
	input.AST = transform.New(&input.DMMF)
	d, _ := json.MarshalIndent(input.AST, "", "  ")
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("AST: %s\n", string(d))
	}
}

func (r *Root) Output(name types.String) dmmf.SchemaField {
	for _, input := range r.DMMF.Schema.OutputObjectTypes.Prisma {
		log.Printf("%s", input.Name)
		if input.Name == "Query" {
			for _, field := range input.Fields {
				if field.Name == name {
					return field
				}
			}
		}
	}
	panic("no such key found: " + name)
}

func (r *Root) Input(name types.String) dmmf.CoreType {
	for _, input := range r.DMMF.Schema.InputObjectTypes.Prisma {
		log.Printf("%s", input.Name)
		if input.Name == name {
			return input
		}
	}
	panic("no such key found: " + name)
}
