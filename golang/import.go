package golang

import (
	"os"
	"path/filepath"
	"strings"

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
	var imports []string
	for _, v := range i {
		imports = append(imports, v.Get())
	}
	return imports
}

// GoImports executes imports missing packages and formats code with goimports command in the given path
func GoImports(projectPath string, commit git.CallbackAddAndCommit) error {
	cmd := getGoImportsCommand(projectPath)
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return err
	}

	var filenames []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(info.Name(), ".go") {
			filenames = append(filenames, strings.Replace(path, projectPath, "", -1))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return commit(filenames, "adding missing imports and format code")
}
