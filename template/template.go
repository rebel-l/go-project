package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-utils/osutils"
)

func CreateFilesWithTemplatePath(pattern string, path string, templateFileConfig map[string]string, data interface{}) ([]string, error) {
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %s", err)
	}

	return CreateFiles(path, tmpl, templateFileConfig, data)
}

func CreateFiles(path string, tmpl *template.Template, templateFileConfig map[string]string, data interface{}) ([]string, error) {
	var fileList []string

	for k, v := range templateFileConfig {
		filename, err := CreateFile(path, tmpl, k, v, data)
		if err != nil {
			return nil, err
		}

		fileList = append(fileList, filename)
	}

	return fileList, nil
}

func CreateFile(path string, tmpl *template.Template, tmplKey string, filename string, data interface{}) (string, error) {
	filename = filepath.Join(path, filename)
	if osutils.FileOrPathExists(filename) {
		return "", nil
	}

	subPath := filepath.Dir(filename)
	if path != subPath {
		if err := osutils.CreateDirectoryIfNotExists(subPath); err != nil {
			return "", err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplKey, data); err != nil {
		return "", fmt.Errorf("failed to write template to file %s: %w", filename, err)
	}

	return filename, nil
}
