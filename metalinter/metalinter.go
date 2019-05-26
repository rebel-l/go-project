package metalinter

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-project/git"
)

// Init initialises the metalinter
func Init(projectPath string, commit git.CallbackAddAndCommit) error {
	pattern := filepath.Join("./metalinter/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(projectPath, "gometalinter.json")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create gometalinter config: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "gometalinter", nil); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	return commit([]string{filename}, "added config for metalinter")
}
