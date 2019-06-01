package readme

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/config"
)

// Init intialises th readme file
func Init(projectPath string, cfg config.Config, commit git.CallbackAddAndCommit) error {
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

	if err = tmpl.ExecuteTemplate(file, "readme", cfg); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}
	return commit([]string{filename}, "added readme")
}
