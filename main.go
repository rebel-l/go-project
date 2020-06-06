package main

import (
	"flag"
	"fmt"

	"github.com/rebel-l/go-project/model"

	"github.com/rebel-l/go-project/docker"

	"github.com/rebel-l/go-project/vagrant"

	"github.com/fatih/color"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/rebel-l/go-project/code"
	"github.com/rebel-l/go-project/description"
	"github.com/rebel-l/go-project/destination"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
	"github.com/rebel-l/go-project/golangci"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/config"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-project/license"
	"github.com/rebel-l/go-project/readme"
	"github.com/rebel-l/go-project/scripts"
	"github.com/rebel-l/go-project/travisci"
)

var (
	pushToRemote *bool
	createModel  *bool
)

func main() {
	// setup parameters
	pushToRemote = flag.Bool("nopush", false, "avoids pushing to remote origin which is helpful for development")
	createModel = flag.Bool("model", false, "creates a new database model for an existing go project")
	flag.Parse()

	// introduction
	fmt.Println()
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("Welcome to Go-Project Tool ...")
	fmt.Println()

	// destination path
	if err := destination.Init(); err != nil {
		print.Error("Init destination path failed", err)
		return
	}
	fmt.Println()

	if *createModel {
		if err := model.Init(destination.Get()); err != nil {
			print.Error("create model failed", err)
		}
		return
	}

	// project kind & description
	kind.Init()
	fmt.Println()
	description.Init()
	fmt.Println()

	// git setup
	git.Setup(destination.Get(), kind.Get())

	// license
	if err := license.Init(destination.Get(), git.GetAuthor().Name, description.Get(), git.AddFilesAndCommit, 0); err != nil {
		print.Error("Init license failed", err)
		return
	}

	// project setup
	setupProject()
	fmt.Println()

	// finish
	print.Info("... Go-Project Tool finished successful!\n")
	print.Warning(
		fmt.Sprintf(
			"Please switch to your project %s and execute './scripts/tools/setup.sh' to install the global tools (Windows users can use GitBash)\n",
			destination.Get(),
		),
	)
	print.Warning("You need to activate this project manually for these third party tools:\nhttps://travis-ci.org\nhttps://codecov.io\n")
}

func setupProject() {
	cfg := config.New(git.GetRemote(), description.Get(), git.GetAuthor())
	vagrant.Prepare(cfg.Project)
	docker.Prepare(cfg.Project)
	fmt.Println()

	total := 10
	if !*pushToRemote {
		total++
	}
	bar := pb.StartNew(total)

	// 1 - main gitignore
	if err := git.CreateIgnore(destination.Get(), git.IgnoreMain, "main gitignore", 1); err != nil {
		print.Error("Create main gitignore failed", err)
		return
	}
	bar.Increment()

	// 2 - golangci
	if err := golangci.Init(destination.Get(), git.AddFilesAndCommit, 2); err != nil {
		print.Error("Create golangci config failed", err)
		return
	}
	bar.Increment()

	// 3 - travis ci
	if err := travisci.Init(destination.Get(), git.AddFilesAndCommit, 3); err != nil {
		print.Error("Create travis file failed", err)
		return
	}
	bar.Increment()

	// 4 - readme
	if err := readme.Init(destination.Get(), cfg, license.Get(), git.AddFilesAndCommit, 4); err != nil {
		print.Error("Create readme failed", err)
		return
	}
	bar.Increment()

	// 5 - go mod
	if err := golang.Init(destination.Get(), cfg.GetPackage(), git.AddFilesAndCommit, 5); err != nil {
		print.Error("Create go mod failed", err)
		return
	}
	bar.Increment()

	git.GetRemote()

	// 6 - vagrant for docker
	if err := vagrant.Setup(destination.Get(), git.AddFilesAndCommit, 6); err != nil {
		print.Error("Create vagrant failed", err)
		return
	}
	bar.Increment()

	// 7 -  docker
	if err := docker.Setup(destination.Get(), git.AddFilesAndCommit, 7); err != nil {
		print.Error("Create docker failed", err)
		return
	}
	bar.Increment()

	// 8 - code
	if err := code.Init(kind.Get(), destination.Get(), cfg, license.Get(), golang.Get, git.AddFilesAndCommit, 8); err != nil {
		print.Error("Creating code base failed", err)
		return
	}
	bar.Increment()

	// 9 - scripts
	if err := scripts.Init(destination.Get(), git.AddFilesAndCommit, git.CreateIgnore, 9); err != nil {
		print.Error("Create scripts failed", err)
		return
	}
	bar.Increment()

	// 10 - run goimports to import missing go packages and format code
	if err := golang.GoImports(destination.Get(), git.AddFilesAndCommit, 10); err != nil {
		print.Error("Formatting code failed", err)
		return
	}
	bar.Increment()

	// 11 - final step: push to remote
	if !*pushToRemote {
		if err := git.Finalize(destination.Get()); err != nil {
			print.Error("Pushing to remote failed", err)
		}
		bar.Increment()
	}

	bar.Finish()
}

/*
other TODO:
1. exit with proper Exit Codes
2. Fix: cyclomatic complexity 12 of function setupProject() is high (> 10) (gocyclo)
3. Ensure templates are compiled in the binary created by `go install`
4. Redesign: Use command pattern
*/
