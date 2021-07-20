package transform

func (r *AST) scalars() []string {
	return []string{
		"String",
		"Bool",
		"Int",
		"Float",
		"DateTime",
	}
}
