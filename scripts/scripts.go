package scripts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rebel-l/go-project/git"
)

// Init initialises the necessary scripts
func Init(projectPath string, commitCallback git.CallbackAddAndCommit, ignoreCallback git.CallbackCreateIgnore) error {
	pattern := filepath.Join("./scripts/tmpl", "*", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	scriptPath := filepath.Join(projectPath, "scripts")
	for _, s := range getScripts() {
		dir := filepath.Join(scriptPath, s.folder)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("faied to create folder %s: %s", dir, err)
		}

		if err = createScript(s.getFilenameWithPath(scriptPath), tmpl, s.getTemplateName()); err != nil {
			return fmt.Errorf("faied to create script %s: %s", s.getFilename(), err)
		}
	}

	if err = commitCallback(getScripts().getFilenames(scriptPath), "added scripts"); err != nil {
		return fmt.Errorf("failed to add and commit scripts: %s", err)
	}

	return createReportsFolder(projectPath, ignoreCallback)
}

func createScript(filename string, tmpl *template.Template, tmplName string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create script file %s: %s", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	return tmpl.ExecuteTemplate(file, tmplName, nil)
}

func getScripts() scripts {
	return scripts{
		{
			folder: "hooks",
			name:   "pre-push",
		},
		{
			folder: "tools",
			name:   "setup",
			suffix: "sh",
		},
		{
			folder: "tools",
			name:   "runChecks",
			suffix: "sh",
		},
	}
}

type scripts []script

func (s scripts) getFilenames(rootPath string) []string {
	var filenames []string
	for _, f := range s {
		filenames = append(filenames, f.getFilenameWithPath(rootPath))
	}
	return filenames
}

type script struct {
	folder string
	name   string
	suffix string
}

func (s script) getFilenameWithPath(rootPath string) string {
	return filepath.Join(rootPath, s.folder, s.getFilename())
}

func (s script) getFilename() string {
	if s.suffix == "" {
		return s.name
	}
	return fmt.Sprintf("%s.%s", s.name, s.suffix)
}

func (s script) getTemplateName() string {
	return strings.Join([]string{s.folder, s.name}, ".")
}
