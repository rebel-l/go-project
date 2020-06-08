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

func (s *sql) Generate(m *model) error {
	destPath := path.Join(s.rootPath, "scripts", "sql")
	fileName := fmt.Sprintf("%s_create_%ss.sql", time.Now().Format("20060102"), strings.ToLower(m.Name))

	tmplName := filepath.Join(".", "model", "tmpl", "sql.tmpl")
	funcMap := template.FuncMap{
		"ToLower":     strings.ToLower,
		"SqliteField": SqliteField,
	}

	tmpl, err := template.New("sql").Funcs(funcMap).ParseFiles(tmplName)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	return s.sqlite(m, tmpl, destPath, fileName)
}

func (s *sql) sqlite(m *model, tmpl *template.Template, destPath, fileName string) error {
	destPath = path.Join(destPath, "sqlite")
	fileName = path.Join(destPath, fileName)
	if err := osutils.CreateDirectoryIfNotExists(destPath); err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create sqlite file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "sqlite", m); err != nil {
		return err
	}

	return nil
}

func SqliteField(f *field) string {
	line := strings.ToLower(f.Name)

	switch f.FieldType {
	case fieldTypeUUID:
		line += " CHAR(36)"
	case fieldTypeString:
		line += fmt.Sprintf(" VARCHAR(%d)", f.MaxLength)
	case fieldTypeBool,
		fieldTypeInt:
		line += " INTEGER"
	case fieldTypeFloat:
		line += " REAL"
	case fieldTypeTime:
		line += " DATETIME"
	}

	if !f.Nullable {
		line += " NOT NULL"
	}

	if f.DefaultValue != "" {
		switch f.FieldType {
		case fieldTypeString:
			line += fmt.Sprintf(" DEFAULT '%s'", f.DefaultValue)
		case fieldTypeBool:
			if strings.ToLower(f.DefaultValue) == "true" {
				line += " DEFAULT 1"
			} else {
				line += " DEFAULT 0"
			}
		default:
			line += fmt.Sprintf(" DEFAULT %s", f.DefaultValue)
		}
	}

	if f.PrimaryKey {
		line += " PRIMARY KEY"
		if f.FieldType == fieldTypeInt {
			line += " AUTOINCREMENT"
		}
	}

	return line
}
