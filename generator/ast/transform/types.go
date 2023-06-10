package transform

func (r *AST) scalars() []string {
	var scalars []string
	for _, item := range r.dmmf.Schema.InputObjectTypes.Prisma {
		for _, field := range item.Fields {
			for _, input := range field.InputTypes {
				if input.Location != "scalar" {
					continue
				}

				name := input.Type.String()

				var exists bool
				for _, s := range scalars {
					// prevent duplicate items
					if s == name {
						exists = true
					}
				}

				if exists {
					continue
				}

				scalars = append(scalars, name)
			}
		}
	}

	return scalars
}
