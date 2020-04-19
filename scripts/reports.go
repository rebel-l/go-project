package scripts

import (
	"os"
	"path/filepath"

	"github.com/rebel-l/go-project/git"
)

func createReportsFolder(projectPath string, createIgnore git.CallbackCreateIgnore, step int) error {
	path := filepath.Join(projectPath, "reports")
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}
	return createIgnore(path, git.IgnoreEmptyFolder, "exclude all files from reports", step)
}
