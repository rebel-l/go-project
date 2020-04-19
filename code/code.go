package code

import (
	"github.com/rebel-l/go-project/code/pkg"
	"github.com/rebel-l/go-project/code/service"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/config"
	"github.com/rebel-l/go-project/license"
)

var goGetCallback golang.CallbackGoGet
var commitCallback git.CallbackAddAndCommit
var commitStep int

// Init creates the code base files.
func Init(
	projectKind string,
	projectPath string,
	cfg config.Config,
	license license.License,
	goGet golang.CallbackGoGet,
	commit git.CallbackAddAndCommit,
	step int) error {

	goGetCallback = goGet
	commitCallback = commit
	commitStep = step

	var packages []string
	var err error
	files := []string{"go.mod"}
	switch projectKind {
	case kind.Package:
		params := pkg.NewParameters(cfg, license)
		err = pkg.Create(projectPath, params, commit, commitStep)
	case kind.Service:
		packages = service.GetPackages().GetNames()
		files = append(files, "go.sum")
		params := service.NewParameters(cfg, license)
		err = service.Create(projectPath, params, commit, commitStep)
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

	return commitCallback(files, "added go packages", commitStep)
}
