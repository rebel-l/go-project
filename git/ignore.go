package git

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rebel-l/go-utils/option"
	"github.com/rebel-l/go-utils/osutils"
)

const (
	// IgnoreMain represents the ignore type for main
	IgnoreMain = "main"
	// IgnoreEmptyFolder represents the ignore type for empty folder
	IgnoreEmptyFolder = "empty_folder"
)

// CallbackCreateIgnore defines the callback to create ignore files
type CallbackCreateIgnore func(path, ignoreType, commitMsg string) error

// CreateIgnore create an git ignore file on the given path and ignore type
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

	return AddFilesAndCommit([]string{filename}, commitMsg)
}

func getPossibleIgnoreTypes() option.Options {
	return option.Options{
		{Key: IgnoreMain},
		{Key: IgnoreEmptyFolder},
	}
}
