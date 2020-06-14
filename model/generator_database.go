package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-utils/osutils"
)

type database struct {
	rootPath string
}

func (d *database) Generate(_ *model) error {
	configPath := path.Join(d.rootPath, "config")
	if osutils.FileOrPathExists(configPath) {
		return nil
	}

	if err := osutils.CreateDirectoryIfNotExists(configPath); err != nil {
		return err
	}

	tmplFile := filepath.Join("./model/tmpl", "database.tmpl")

	tmpl, err := template.New("store").ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	for _, tmplID := range getDatabaseTemplateIdentifiers() {
		if err := d.database(tmpl, configPath, tmplID); err != nil {
			return err
		}
	}

	return nil
}

func (d *database) database(tmpl *template.Template, path, tmplID string) error {
	fileName := filepath.Join(path, tmplID)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create store file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, nil); err != nil {
		return err
	}

	return nil
}

func getDatabaseTemplateIdentifiers() []string {
	return []string{
		"package.go",
		"database.go",
		"database_test.go",
	}
}
