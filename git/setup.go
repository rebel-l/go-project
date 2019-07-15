package git

import (
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/osutils"
)

const (
	developBranch = "refs/heads/develop"
)

var repo *git.Repository
var author *object.Signature
var errMsg = color.New(color.FgRed, color.Italic)
var rootPath string
var remote string

// GetAuthor returns the author entered by user. Is nil as long Setup() not called
func GetAuthor() *object.Signature {
	return author
}

// GetRemote returns the remote url to git repository
func GetRemote() string {
	return remote
}

// Setup ensures that git repo is created and remote origin is set
func Setup(projectPath string) {
	rootPath = projectPath
	if !open(projectPath) {
		if err := createRepo(projectPath); err != nil {
			_, _ = errMsg.Printf("Failed to init git repo: %s\n", err) // Rework error message handling, use lib
			return
		}
	}

	author = askForAuthor()
	remote = askForRemote()

	ok, err := hasRemote()
	if err != nil {
		_, _ = errMsg.Printf("\n", err)
		return
	}

	if !ok {
		// TODO: do as last step together with push
		if err = createRemote(); err != nil {
			_, _ = errMsg.Printf("Failed to set remote origin on repo: %s\n", err)
			return
		}
	}

	//if err = createBranch(); err != nil {
	//	_, _ = errMsg.Printf("Failed to set branch on repo: %s\n", err)
	//	return
	//}

	//if err = checkoutBranch(); err != nil {
	//	_, _ = errMsg.Printf("Failed to create and checkout branch on repo: %s\n", err)
	//	return
	//}

	//list, err := repo.Branches()
	//if err != nil {
	//	_, _ = errMsg.Printf("Failed to retrieve branches: %s\n", err)
	//}
	//defer list.Close()
	//err = list.ForEach(func(reference *plumbing.Reference) error {
	//	fmt.Println(reference.String())
	//	fmt.Println(reference.Name())
	//	fmt.Println(reference.Type())
	//	fmt.Println(reference.Hash())
	//	return nil
	//})
	//if err != nil {
	//	_, _ = errMsg.Printf("Failed to iterate branches: %s\n", err)
	//}
	//fmt.Println("TEST")
}

func Finalize() {
	// TODO: push to remote
}

func open(path string) bool {
	if !osutils.FileOrPathExists(filepath.Join(path, git.GitDirName)) {
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
	if remote == "" {
		return nil
	}

	_, err := repo.CreateRemote(&config.RemoteConfig{Name: git.DefaultRemoteName, URLs: []string{remote}})
	return err

	// TODO: pull for existing repo
	//workingTree, err := repo.Worktree()
	//if err != nil {
	//	return err
	//}
	//
	//return workingTree.Pull(&git.PullOptions{RemoteName: git.DefaultRemoteName})
}

func askForRemote() string {
	t := prompt.Input("Enter the remote origin of your branch (leave empty to add it later by yourself): ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})

	return strings.TrimSpace(strings.ToLower(t))
}

func askForAuthor() *object.Signature {
	name := prompt.Input("Enter your git username (used as author and for license): ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})

	email := prompt.Input("Enter your git email (used as author): ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})

	return &object.Signature{Name: name, Email: email}
}

func createBranch() error {
	head, err := repo.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference(developBranch, head.Hash())

	return repo.Storer.SetReference(ref)
}

func checkoutBranch() error {
	workingTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return workingTree.Checkout(&git.CheckoutOptions{
		Branch: developBranch,
		Create: true,
	})
}
