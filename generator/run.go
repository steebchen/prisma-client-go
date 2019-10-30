package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"strings"
	"text/template"
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
		return fmt.Errorf("could not get main template asset: %w", err)
	}

	templateDir := pkg.Dir + "/generator/templates/*.gotpl"
	templates, err := template.ParseGlob(templateDir)
	if err != nil {
		return fmt.Errorf("could not parse go templates dir %s: %w", templateDir, err)
	}

	// Run header template first
	header := templates.Lookup("_header.gotpl")

	err = header.Execute(&buf, input)
	if err != nil {
		return fmt.Errorf("could not write header template %s: %w", input.Generator.Output, err)
	}

	// Then process all remaining templates
	for _, tpl := range templates.Templates() {
		if strings.Contains(tpl.Name(), "_") {
			continue
		}
		buf.Write([]byte(fmt.Sprintf("// --- template %s ---\n", tpl.Name())))
		err = tpl.Execute(&buf, input)
		if err != nil {
			return fmt.Errorf("could not write template file %s: %w", input.Generator.Output, err)
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not format source: %s", err)
	}

	err = ioutil.WriteFile(input.Generator.Output, formatted, 0644)
	if err != nil {
		return fmt.Errorf("could not write template data to file writer %s: %w", input.Generator.Output, err)
	}

	return nil
}
