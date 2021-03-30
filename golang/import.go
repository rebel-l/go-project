package golang

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/rebel-l/go-project/git"
)

// Import represents the information of a go import
type Import struct {
	Name  string
	Alias string
}

// Get returns the import command for go. If an alias is set, it returns it with alias, otherwise only the name of the package
func (i Import) Get() string {
	if i.Alias == "" {
		return "\"" + i.Name + "\""
	}

	return i.Alias + " \"" + i.Name + "\""
}

// Imports is a collection of Import
type Imports []Import

// GetNames returns a list of package names
func (i Imports) GetNames() []string {
	var names []string
	for _, v := range i {
		names = append(names, v.Name)
	}
	return names
}

// Get returns the list of import commands for go as it is provided by Import.Get()
func (i Imports) Get() []string {
	var imp []string
	for _, v := range i {
		imp = append(imp, v.Get())
	}
	return imp
}

// GoImports executes imports missing packages and formats code with goimports command in the given path
func GoImports(projectPath string, commit git.CallbackAddAndCommit, step int) error {
	var filenames []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(info.Name(), ".go") && !strings.Contains(info.Name(), ".json") {
			src, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			p, err := git.GetPackage(projectPath)
			if err != nil {
				return err
			}

			imports.LocalPrefix = p
			res, err := imports.Process(path, src, nil)
			if err != nil {
				return err
			}

			if !bytes.Equal(src, res) {
				if err := ioutil.WriteFile(path, res, 0644); err != nil {
					return err
				}
			}

			filenames = append(filenames, strings.Replace(path, projectPath, "", -1))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err = ModTidy(projectPath); err != nil {
		return err
	}

	return commit(filenames, "adding missing imports and format code", step)
}
