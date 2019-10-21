package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func addDefaults(input *Root) {
	if input.Generator.Config.Package == "" {
		input.Generator.Config.Package = "main"
	}
}

// Run invokes the generator, which builds the templates and writes to the specified output file.
func Run(input Root) error {
	addDefaults(&input)

	var buf bytes.Buffer

	ctx := build.Default
	pkg, err := ctx.Import("github.com/prisma/photongo", ".", build.FindOnly)
	if err != nil {
		return errors.Wrap(err, "could not get main template asset")
	}

	templateDir := pkg.Dir + "/generator/templates/*.gotpl"
	templates, err := template.ParseGlob(templateDir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not parse go templates dir %s", templateDir))
	}

	// Run header template first
	header := templates.Lookup("_header.gotpl")

	err = header.Execute(&buf, input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write header template %s", input.Generator.Output))
	}

	// Then process all remaining templates
	for _, tpl := range templates.Templates() {
		if strings.Contains(tpl.Name(), "_") {
			continue
		}
		buf.Write([]byte(fmt.Sprintf("// --- template %s ---\n", tpl.Name())))
		err = tpl.Execute(&buf, input)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not write template file %s", input.Generator.Output))
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not format source"))
	}

	err = ioutil.WriteFile(input.Generator.Output, formatted, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write template data to file writer %s", input.Generator.Output))
	}

	return nil
}
