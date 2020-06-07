package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rebel-l/go-utils/osutils"
)

type store struct {
	rootPath string
}

func (s *store) Generate(m *model) error {
	name := strings.ToLower(m.Name)

	destPath := path.Join(s.rootPath, name, name+"store")
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return err
	}

	tmplFile := filepath.Join("./model/tmpl", "store.tmpl")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		//"StoreField":
	}

	tmpl, err := template.New("store").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	for _, tmplID := range getStoreTemplateIdentifiers() {
		if err := s.store(m, tmpl, destPath, tmplID); err != nil {
			return err
		}
	}

	return nil
}

func (s store) store(m *model, tmpl *template.Template, path, tmplID string) error {
	fileName := filepath.Join(path, strings.Replace(tmplID, "store", strings.ToLower(m.Name), -1))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create store file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, m); err != nil {
		return err
	}

	return nil
}

func getStoreTemplateIdentifiers() []string {
	return []string{
		"package.go",
		"store.go",
		"store_test.go",
	}
}
