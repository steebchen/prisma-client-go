package builder

func TransformEquals(fields []Field) []Field {
	for i, field := range fields {
		if field.Fields != nil {
			for _, inner := range field.Fields {
				if inner.Name == "equals" {
					fields[i].Value = inner.Value
					fields[i].Fields = nil
				}
			}
		}
	}
	return fields
}
