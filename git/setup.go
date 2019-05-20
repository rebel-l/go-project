package git

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

const remoteName = "origin"

var repo *git.Repository
var errMsg = color.New(color.FgRed, color.Italic)

// Setup ensures that git repo is created and remote origin is set
func Setup(projectPath string) {
	if !open(projectPath) {
		if err := createRepo(projectPath); err != nil {
			_, _ = errMsg.Printf("Failed to init git repo: %s\n", err)
			return
		}
	}

	ok, err := hasRemote()
	if err != nil {
		_, _ = errMsg.Printf("\n", err)
		return
	}

	if !ok {
		if err = createRemote(); err != nil {
			_, _ = errMsg.Printf("Failed to set remote origin on repo: %s\n", err)
			return
		}
	}
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

func hasRemote() (bool, error) {
	rem, err := repo.Remotes()
	if err != nil {
		return false, err
	}

	return len(rem) > 0, nil
}

func createRemote() error {
	remote := askForRemote()
	if remote == "" {
		return nil
	}

	_, err := repo.CreateRemote(&config.RemoteConfig{Name: remoteName, URLs: []string{remote}})
	if err != nil {
		return err
	}

	return repo.Fetch(&git.FetchOptions{RemoteName: remoteName})
}

func askForRemote() string {
	t := prompt.Input("Enter the remote origin of your branch (leave empty to add it later by yourself): ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})

	return strings.TrimSpace(strings.ToLower(t))
}
