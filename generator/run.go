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
	fmt.Printf("package: %s\n", input.Generator.Config.Package)
	fmt.Printf("file: %s\n", input.Generator.Output)
	input.Generator.Config.Package = "main"
}

func Run(input Root) error {
	addDefaults(&input)

	asset, err := templates.Asset("generator/templates/main.gotpl")
	if err != nil {
		return errors.Wrap(err, "could not get main template asset")
	}

	tpl, err := template.New("main").Parse(string(asset))
	if err != nil {
		return errors.Wrap(err, "could not parse templates")
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write template data to file writer %s", input.Generator.Output))
	}

	err = ioutil.WriteFile(input.Generator.Output, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not write template data to file writer %s", input.Generator.Output))
	}

	err = exec.Command("go", "fmt", input.Generator.Output).Run()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("could not format file with go fmt %s", input.Generator.Output))
	}

	return nil
}
