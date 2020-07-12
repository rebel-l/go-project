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

type modelGen struct {
	rootPath string
}

func (g *modelGen) Generate(m *model) error {
	name := strings.ToLower(m.Name)

	destPath := path.Join(g.rootPath, name, name+"model")
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return err
	}

	tmplFile := filepath.Join("./model/tmpl", "model.tmpl")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New("model").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	for _, tmplID := range getModelTemplateIdentifiers() {
		if err := g.model(m, tmpl, destPath, tmplID); err != nil {
			return err
		}
	}

	return nil
}

func (g *modelGen) model(m *model, tmpl *template.Template, path, tmplID string) error {
	fileName := filepath.Join(path, strings.Replace(tmplID, "model", strings.ToLower(m.Name), -1))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create model file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, m); err != nil {
		return err
	}

	return nil
}

func getModelTemplateIdentifiers() []string {
	return []string{
		"package.go",
		"model.go",
		"model_test.go",
	}
}
