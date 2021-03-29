package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-utils/osutils"
)

type config struct {
	rootPath string
}

func (d *config) Generate(m *model) ([]string, error) {
	configPath := path.Join(d.rootPath, "config")
	if osutils.FileOrPathExists(configPath) {
		return nil, nil
	}

	if err := osutils.CreateDirectoryIfNotExists(configPath); err != nil {
		return nil, err
	}

	tmplFile := filepath.Join("./model/tmpl", "config.tmpl")

	tmpl, err := template.New("config").ParseFiles(tmplFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %s", err)
	}

	var files []string
	for _, tmplID := range getDatabaseTemplateIdentifiers() {
		f, err := d.config(tmpl, m, configPath, tmplID)
		if err != nil {
			return nil, err
		}

		if f != "" {
			files = append(files, f)
		}
	}

	return files, nil
}

func (d *config) config(tmpl *template.Template, m *model, path, tmplID string) (string, error) {
	fileName := filepath.Join(path, tmplID)
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create config file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, m); err != nil {
		return "", err
	}

	return fileName, nil
}

func getDatabaseTemplateIdentifiers() []string {
	return []string{
		"package.go",
		"database.go",
		"database_test.go",
	}
}
