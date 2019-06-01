package golang

import (
	"github.com/rebel-l/go-project/git"
)

// Init initialises go modules
func Init(projectPath string, packageName string, commit git.CallbackAddAndCommit) error {
	cmd := getGoModCommand(packageName)
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return err
	}
	return commit([]string{"go.mod"}, "init go mod")
}

// GoGet defines the callback for go get
type GoGet func(projectPath string, packageName string) error

// Get executes go get on for given package name
func Get(projectPath string, packageName string) error {
	cmd := getGoGetCommand(packageName)
	cmd.Dir = projectPath
	return cmd.Run()
}
