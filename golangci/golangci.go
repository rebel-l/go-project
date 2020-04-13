package golangci

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/git"
)

// Init initialises the golangci
func Init(projectPath string, commit git.CallbackAddAndCommit, step int) error {
	pattern := filepath.Join("./golangci/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(projectPath, ".golangci.json")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create golangci config: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "golangci", nil); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	return commit([]string{filename}, "added config for golangci", step)
}
