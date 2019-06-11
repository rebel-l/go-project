package code

import (
	"github.com/rebel-l/go-project/code/pkg"
	"github.com/rebel-l/go-project/code/service"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/config"
)

var goGetCallback golang.CallbackGoGet
var commitCallback git.CallbackAddAndCommit

// Init creates the code base files.
func Init(projectKind string, projectPath string, cfg config.Config, goGet golang.CallbackGoGet, commit git.CallbackAddAndCommit) error {
	goGetCallback = goGet
	commitCallback = commit

	var packages []string
	var err error
	files := []string{"go.mod"}
	switch projectKind {
	case kind.Package:
		err = pkg.Create(projectPath, cfg, commit)
	case kind.Service:
		packages = service.GetPackages().GetNames()
		files = append(files, "go.sum")
		params := service.NewParameters(cfg)
		err = service.Create(projectPath, params, commit)
	}

	if err != nil {
		return err
	}

	return addPackages(packages, projectPath, files)
}

func addPackages(packages []string, projectPath string, files []string) error {
	for _, p := range packages {
		err := goGetCallback(projectPath, p)
		if err != nil {
			return err
		}
	}

	return commitCallback(files, "added go packages")
}
