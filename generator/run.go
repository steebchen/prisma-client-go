package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

func addDefaults(input *Root) {
	fmt.Printf("package: %s\n", input.Generator.Config.Package)
	fmt.Printf("file: %s\n", input.Generator.Output)
	input.Generator.Config.Package = "main"
}

func Run(input Root) error {
	addDefaults(&input)

	exec, err := os.Executable()
	dir := filepath.Dir(exec)
	if err != nil {
		return errors.Wrap(err, "could not get executable")
	}

	tpl, err := template.ParseGlob(dir + "/generator/templates/*.gotpl")
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

	return nil
}
