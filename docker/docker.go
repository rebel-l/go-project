package docker

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/rebel-l/go-project/dialog"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/template"

	"github.com/rebel-l/go-project/kind"
)

const (
	templateKeyDocker        = "Dockerfile"
	templateKeyScriptRestart = "scriptRestart"
)

var (
	params *Docker
)

type Docker struct {
	ServiceName string
}

func Prepare(serviceName string) {
	if kind.Get() != kind.Service || !dialog.Confirmation("Add Docker to this project?") {
		return
	}

	params = &Docker{ServiceName: serviceName}
}

func Setup(destination string, commit git.CallbackAddAndCommit, step int) error {
	if params == nil {
		return nil
	}

	fileConfig := map[string]string{
		templateKeyDocker:        "Dockerfile",
		templateKeyScriptRestart: path.Join("scripts", "tools", "restartService.sh"),
	}

	pattern := filepath.Join("./docker/tmpl", "*.tmpl")
	filenames, err := template.CreateFilesWithTemplatePath(pattern, destination, fileConfig, params)
	if err != nil {
		return fmt.Errorf("docker setup failed: %w", err)
	}

	return commit(filenames, "setup docker", step)
}
