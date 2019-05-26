package git

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/rebel-l/go-utils/osutils"

	"github.com/rebel-l/go-utils/option"
)

const (
	IgnoreMain        = "main"
	IgnoreEmptyFolder = "empty_folder"
)

func CreateIgnore(path, ignoreType, commitMsg string) error {
	if !osutils.FileOrPathExists(path) {
		return fmt.Errorf("path %s doesn't exist", path)
	}

	if !getPossibleIgnoreTypes().IsValidOption(ignoreType) {
		return fmt.Errorf("%s is not a valid type", ignoreType)
	}

	pattern := filepath.Join("./git/tmpl/ignore", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(path, ".gitignore")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, ignoreType, nil); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	return AddFilesAndCommit([]string{".gitignore"}, commitMsg)
}

func getPossibleIgnoreTypes() option.Options {
	return option.Options{
		{Key: IgnoreMain},
		{Key: IgnoreEmptyFolder},
	}
}
