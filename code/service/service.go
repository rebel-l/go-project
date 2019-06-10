package service

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/config"
)

// Parameters defines parameters used for the go templates
type Parameters struct {
	LicensePrefix string
	Packages      []string
}

// NewParameters returns a new struct of Parameters prefilled by a config and the definition of packages
func NewParameters(cfg config.Config) Parameters {
	return Parameters{
		LicensePrefix: cfg.LicensePrefix,
		Packages:      GetPackages(),
	}
}

// Create the basic files for a service
func Create(projectPath string, params Parameters, commit git.CallbackAddAndCommit) error {
	filename := filepath.Join(projectPath, "main.go")
	pattern := filepath.Join("./code/service/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create service main file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "main", params); err != nil {
		return err
	}
	return commit([]string{filename}, "added main go file for service")
}

/*
TODO:
1. logrus
2. mux
3. ping endpoint
4. flags with port parameter
5. test file for ping endpoint
6. swagger definition
7. later: auth client - permission request
*/
