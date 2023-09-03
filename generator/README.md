# Generator

This package handles the actual generation of the Go client files. It handles copying engines and converting templates
with a given AST to a Go client ORM file.

Note that there is a lot of special logic around "DMMF" as the design was initially intended for JavaScript, but not for
type-safe languages like Go.
