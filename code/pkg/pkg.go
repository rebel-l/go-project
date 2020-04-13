package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/license"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/config"
)

// Parameters defines parameters used for the go templates
type Parameters struct {
	Config  config.Config
	License license.License
}

// NewParameters returns a new struct of Parameters prefilled by a config and license
func NewParameters(cfg config.Config, license license.License) Parameters {
	return Parameters{
		Config:  cfg,
		License: license,
	}
}

// Create creates the basic files for a package
func Create(projectPath string, params Parameters, commit git.CallbackAddAndCommit, step int) error {
	filename := filepath.Join(projectPath, fmt.Sprintf("%s.go", params.Config.Project))
	pattern := filepath.Join("./code/pkg/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create package main file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "pkg", params); err != nil {
		return err
	}
	return commit([]string{filename}, "added main go file for package", step)
}
