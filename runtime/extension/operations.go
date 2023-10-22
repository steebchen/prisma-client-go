package extension

import "strings"

var Operations = []Operation{
	CreateOne,
	DeleteMany,
	DeleteOne,
	FindFirst,
	FindMany,
	FindUnique,
	UpdateMany,
	UpdateOne,
	UpsertOne,
}

type Operation string

const (
	CreateOne  Operation = "CreateOne"
	DeleteMany Operation = "DeleteMany"
	DeleteOne  Operation = "DeleteOne"
	FindFirst  Operation = "FindFirst"
	FindMany   Operation = "FindMany"
	FindUnique Operation = "FindUnique"
	UpdateMany Operation = "UpdateMany"
	UpdateOne  Operation = "UpdateOne"
	UpsertOne  Operation = "UpsertOne"
)

func (r Operation) IsMany() bool {
	return strings.HasSuffix(string(r), "Many")
}

func (r Operation) Param() string {
	str := ""
	if !r.IsMany() {
		str += "Unique"
	}
	str += "WhereParam"
	return str
}
