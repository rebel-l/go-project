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

func (b *bootstrap) Generate(m *model) error {
	bootstrapPath := path.Join(b.rootPath, "bootstrap")
	if osutils.FileOrPathExists(bootstrapPath) {
		return nil
	}

	if err := osutils.CreateDirectoryIfNotExists(bootstrapPath); err != nil {
		return err
	}

	tmplFile := filepath.Join("./model/tmpl", "bootstrap.tmpl")

	tmpl, err := template.New("bootstrap").ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	for _, tmplID := range getDatabaseTemplateIdentifiers() {
		if err := b.bootstrap(tmpl, m, bootstrapPath, tmplID); err != nil {
			return err
		}
	}

	return nil
}

func (b *bootstrap) bootstrap(tmpl *template.Template, m *model, path, tmplID string) error {
	fileName := filepath.Join(path, tmplID)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create bootstrap file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, m); err != nil {
		return err
	}

	return nil
}
