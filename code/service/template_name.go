package service

import (
	"path/filepath"
	"strings"
)

type templateName struct {
	Name          string
	FileExtension string
}

func (t templateName) toFilenName() string {
	name := strings.Replace(t.Name, ".", string(filepath.Separator), -1)
	if t.FileExtension != "" {
		name += "." + t.FileExtension
	}
	return name
}
