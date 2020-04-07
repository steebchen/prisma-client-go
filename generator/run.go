// Package generator acts as a prisma generator
package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
)

func addDefaults(input *Root) {
	if input.Generator.Config.Package == "" {
		input.Generator.Config.Package = "main"
	}
}

// Run invokes the generator, which builds the templates and writes to the specified output file.
func Run(input *Root) error {
	addDefaults(input)

	// copy the query engine to the local repository path
	for _, name := range input.Generator.BinaryTargets {
		if name == "native" {
			name = platform.BinaryPlatformName()
		}

		// first, ensure they are actually downloaded
		if err := binaries.FetchEngine(binaries.GlobalTempDir(), "query-engine", name); err != nil {
			return fmt.Errorf("failed fetching binaries: %w", err)
		}

		qe := "prisma-query-engine-" + name

		dir := path.Join(binaries.GlobalTempDir(), qe)

		input, err := ioutil.ReadFile(dir)
		if err != nil {
			return fmt.Errorf("could not read file %s: %w", name, err)
		}

		dest := "./" + qe
		err = ioutil.WriteFile(dest, input, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not write file to %s: %w", dest, err)
		}
	}

	var buf bytes.Buffer

	ctx := build.Default
	pkg, err := ctx.Import("github.com/prisma/prisma-client-go", ".", build.FindOnly)
	if err != nil {
		return fmt.Errorf("could not get main template asset: %w", err)
	}

	var templates []*template.Template

	templateDir := pkg.Dir + "/generator/templates"
	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gotpl") {
			tpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}
			templates = append(templates, tpl.Templates()...)
		}

		return err
	})

	if err != nil {
		return fmt.Errorf("could not walk dir %s: %w", templateDir, err)
	}

	// Run header template first
	header, err := template.ParseFiles(templateDir + "/_header.gotpl")
	if err != nil {
		return fmt.Errorf("could not find header template %s: %w", templateDir, err)
	}

	if err := header.Execute(&buf, input); err != nil {
		return fmt.Errorf("could not write header template: %w", err)
	}

	// Then process all remaining templates
	for _, tpl := range templates {
		if strings.Contains(tpl.Name(), "_") {
			continue
		}

		buf.Write([]byte(fmt.Sprintf("// --- template %s ---\n", tpl.Name())))

		err = tpl.Execute(&buf, input)

		if err != nil {
			return fmt.Errorf("could not write template file %s: %w", tpl.Name(), err)
		}

		_, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("could not format source %s from file %s: %w", buf.String(), tpl.Name(), err)
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not format final source: %w", err)
	}

	outputDir := filepath.Dir(input.Generator.Output)
	if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("could not run MkdirAll on path %s: %w", input.Generator.Output, err)
	}

	if err = ioutil.WriteFile(input.Generator.Output, formatted, 0644); err != nil {
		return fmt.Errorf("could not write template data to file writer %s: %w", input.Generator.Output, err)
	}

	return nil
}
