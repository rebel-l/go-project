package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/rebel-l/go-utils/osutils"
)

type sql struct {
	rootPath string
}

func (s *sql) Generate(m *model) ([]string, error) {
	destPath := path.Join(s.rootPath, "scripts", "sql")
	fileName := fmt.Sprintf("%s_create_%ss.sql", time.Now().Format("20060102_150405"), strings.ToLower(m.Name))

	tmplName := filepath.Join(".", "model", "tmpl", "sql.tmpl")
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New("sql").Funcs(funcMap).ParseFiles(tmplName)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	f, err := s.sqlite(m, tmpl, destPath, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sql files: %w", err)
	}

	if f != "" {
		return []string{f}, nil
	}

	return nil, nil
}

func (s *sql) sqlite(m *model, tmpl *template.Template, destPath, fileName string) (string, error) {
	destPath = path.Join(destPath, "sqlite")
	fileName = path.Join(destPath, fileName)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return "", err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create sqlite file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "sqlite", m); err != nil {
		return "", err
	}

	return fileName, nil
}
