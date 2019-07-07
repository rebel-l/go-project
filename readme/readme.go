package readme

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

// Init intialises th readme file
func Init(projectPath string, cfg config.Config, license license.License, commit git.CallbackAddAndCommit) error {
	pattern := filepath.Join("./readme/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(projectPath, "README.md")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create readme file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	params := NewParameters(cfg, license)
	if err = tmpl.ExecuteTemplate(file, "readme", params); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}
	return commit([]string{filename}, "added readme")
}
