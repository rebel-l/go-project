package git

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
	"gopkg.in/src-d/go-git.v4"
)

var repo *git.Repository
var errMsg = color.New(color.FgRed, color.Italic)

// Setup ensures that git repo is created and remote origin is set
func Setup(projectPath string) {
	if !open(projectPath) {
		if err := createRepo(projectPath); err != nil {
			_, _ = errMsg.Printf("Failed to init git repo: %s", err)
			return
		}
	}

	fmt.Println(repo.Head())
	// TODO: check remote origin
}

func open(path string) bool {
	if !osutils.FileOrPathExists(path + "/" + git.GitDirName) {
		return false
	}

	var err error
	repo, err = git.PlainOpen(path)
	if err != nil {
		return false
	}

	return true
}

func createRepo(path string) error {
	var err error
	repo, err = git.PlainInit(path, false)
	return err
}
