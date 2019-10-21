package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
)

func addDefaults(input *Root) {
	if input.Generator.Config.Package == "" {
		input.Generator.Config.Package = "main"
	}
}

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

	for _, tpl := range templates.Templates() {
		buf.Write([]byte(fmt.Sprintf("// --- template %s ---\n", tpl.Name())))
		err = tpl.Execute(&buf, input)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not write template data to file writer %s", input.Generator.Output))
		}
	}

	if err := ioutil.WriteFile(input.Generator.Output, buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write template data to file writer %s", input.Generator.Output))
	}

	if err := exec.Command("go", "fmt", input.Generator.Output).Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not format file with go fmt %s", input.Generator.Output))
	}

	return nil
}
