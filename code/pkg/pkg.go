package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/config"
)

// Create creates the basic files for a package
func Create(projectPath string, cfg config.Config, commit git.CallbackAddAndCommit) error {
	filename := filepath.Join(projectPath, fmt.Sprintf("%s.go", cfg.Project))
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

	if err = tmpl.ExecuteTemplate(file, "pkg", cfg); err != nil {
		return err
	}
	return commit([]string{filename}, "added main go file for package")
}
