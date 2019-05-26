package travisci

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-project/git"
)

// Init initialises travis ci
func Init(projectPath string, commit git.CallbackAddAndCommit) error {
	pattern := filepath.Join("./travisci/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(projectPath, ".travis.yml")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create travis file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "travis", nil); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	return commit([]string{filename}, "setup travis")
}
