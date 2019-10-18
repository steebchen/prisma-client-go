package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"

	"github.com/prisma/photongo/generator/templates"
)

func addDefaults(input *Root) {
	if input.Generator.Config.Package == "" {
		input.Generator.Config.Package = "main"
	}
}

func Run(input Root) error {
	addDefaults(&input)

	var buf bytes.Buffer

	for _, name := range templates.AssetNames() {
		asset, err := templates.Asset(name)
		if err != nil {
			return errors.Wrap(err, "could not get main template asset")
		}

		tpl, err := template.New(name).Parse(string(asset))
		if err != nil {
			return errors.Wrap(err, "could not parse templates")
		}

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
