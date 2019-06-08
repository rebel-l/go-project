package code

import (
	"github.com/rebel-l/go-project/code/service"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
	"github.com/rebel-l/go-project/kind"
)

var goGetCallback golang.CallbackGoGet
var commitCallback git.CallbackAddAndCommit

func Init(projectKind string, projectPath string, goGet golang.CallbackGoGet, commit git.CallbackAddAndCommit) error {
	goGetCallback = goGet
	commitCallback = commit

	var packages []string
	switch projectKind {
	case kind.Package:
		// TODO
	case kind.Service:
		packages = service.GetPackages()
	}
	return addPackages(packages, projectPath)
}

func addPackages(packages []string, projectPath string) error {
	for _, p := range packages {
		err := goGetCallback(projectPath, p)
		if err != nil {
			return err
		}
	}

	return commitCallback([]string{"go.sum", "go.mod"}, "added go packages")
}
