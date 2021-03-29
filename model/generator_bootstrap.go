package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-utils/osutils"
)

type bootstrap struct {
	rootPath string
}

func (b *bootstrap) Generate(m *model) ([]string, error) {
	bootstrapPath := path.Join(b.rootPath, "bootstrap")
	if osutils.FileOrPathExists(bootstrapPath) {
		return nil, nil
	}

	if err := osutils.CreateDirectoryIfNotExists(bootstrapPath); err != nil {
		return nil, err
	}

	tmplFile := filepath.Join("./model/tmpl", "bootstrap.tmpl")

	tmpl, err := template.New("bootstrap").ParseFiles(tmplFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %s", err)
	}

	var files []string
	for _, tmplID := range getDatabaseTemplateIdentifiers() {
		f, err := b.bootstrap(tmpl, m, bootstrapPath, tmplID)
		if err != nil {
			return nil, err
		}

		if f != "" {
			files = append(files, f)
		}
	}

	return files, nil
}

func (b *bootstrap) bootstrap(tmpl *template.Template, m *model, path, tmplID string) (string, error) {
	fileName := filepath.Join(path, tmplID)
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create bootstrap file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, m); err != nil {
		return "", err
	}

	return fileName, nil
}
