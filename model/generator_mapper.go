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

type mapper struct {
	rootPath string
}

func (m *mapper) Generate(mo *model) ([]string, error) {
	name := strings.ToLower(mo.Name)

	destPath := path.Join(m.rootPath, name, name+"mapper")
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return nil, err
	}

	tmplFile := filepath.Join("./model/tmpl", "mapper.tmpl")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New("mapper").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %s", err)
	}

	var files []string
	for _, tmplID := range getMapperTemplateIdentifiers() {
		f, err := m.mapper(mo, tmpl, destPath, tmplID)
		if err != nil {
			return nil, err
		}

		if f != "" {
			files = append(files, f)
		}
	}

	return files, nil
}

func (m *mapper) mapper(mo *model, tmpl *template.Template, path, tmplID string) (string, error) {
	fileName := filepath.Join(path, strings.Replace(tmplID, "mapper", strings.ToLower(mo.Name), -1))
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create mapper file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplID, mo); err != nil {
		return "", err
	}

	return fileName, nil
}

func getMapperTemplateIdentifiers() []string {
	return []string{
		"package.go",
		"mapper.go",
		"mapper_test.go",
	}
}
