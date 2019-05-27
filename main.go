package main

import (
	"fmt"

	"github.com/rebel-l/go-project/readme"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/destination"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-project/license"
	"github.com/rebel-l/go-project/metalinter"
	"github.com/rebel-l/go-project/scripts"
	"github.com/rebel-l/go-project/travisci"
)

func main() {
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

	// project kind
	kind.Init()
	fmt.Println()

	// git setup
	git.Setup(destination.Get())

	// license
	if err := license.Init(destination.Get(), git.GetAuthor().Name, git.AddFilesAndCommit); err != nil {
		print.Error("Init license failed", err)
		return
	}
	fmt.Println()

	// main gitignore
	if err := git.CreateIgnore(destination.Get(), git.IgnoreMain, "main gitignore"); err != nil {
		print.Error("Create main gitignore failed", err)
		return
	}

	// metalinter
	if err := metalinter.Init(destination.Get(), git.AddFilesAndCommit); err != nil {
		print.Error("Create metalinter config failed", err)
		return
	}

	// travis ci
	if err := travisci.Init(destination.Get(), git.AddFilesAndCommit); err != nil {
		print.Error("Create travis file failed", err)
		return
	}

	// readme
	if err := readme.Init(destination.Get(), git.GetRemote(), git.AddFilesAndCommit); err != nil {
		print.Error("Create readme failed", err)
		return
	}

	// go mod
	// TODO

	// vagrant for docker
	// TODO

	// docker
	// TODO

	// code
	// TODO: package & service

	// scripts
	if err := scripts.Init(destination.Get(), git.AddFilesAndCommit, git.CreateIgnore); err != nil {
		print.Error("Create scripts failed", err)
		return
	}

	// finish
	print.Info("... Go-Project Tool finished successful!\n")
	print.Warning(
		fmt.Sprintf(
			"Please switch to your project %s and execute './scripts/tools/setup.sh' to install the global tools (Windows users can use GitBash)\n",
			destination.Get(),
		),
	)
}
