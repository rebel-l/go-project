package main

import (
	"fmt"

	"github.com/fatih/color"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/rebel-l/go-project/code"
	"github.com/rebel-l/go-project/description"
	"github.com/rebel-l/go-project/destination"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/config"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-project/license"
	"github.com/rebel-l/go-project/metalinter"
	"github.com/rebel-l/go-project/readme"
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

	// project kind & description
	kind.Init()
	fmt.Println()
	description.Init()
	fmt.Println()

	// git setup
	git.Setup(destination.Get())

	// license
	if err := license.Init(destination.Get(), git.GetAuthor().Name, description.Get(), git.AddFilesAndCommit); err != nil {
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
}

func setupProject() {
	cfg := config.New(git.GetRemote(), license.Get(), license.GetPrefix(), description.Get())
	fmt.Println()

	bar := pb.StartNew(10)
	// main gitignore
	if err := git.CreateIgnore(destination.Get(), git.IgnoreMain, "main gitignore"); err != nil {
		print.Error("Create main gitignore failed", err)
		return
	}
	bar.Increment()

	// metalinter
	if err := metalinter.Init(destination.Get(), git.AddFilesAndCommit); err != nil {
		print.Error("Create metalinter config failed", err)
		return
	}
	bar.Increment()

	// travis ci
	if err := travisci.Init(destination.Get(), git.AddFilesAndCommit); err != nil {
		print.Error("Create travis file failed", err)
		return
	}
	bar.Increment()

	// readme
	if err := readme.Init(destination.Get(), cfg, git.AddFilesAndCommit); err != nil {
		print.Error("Create readme failed", err)
		return
	}
	bar.Increment()

	// go mod
	if err := golang.Init(destination.Get(), cfg.GetPackage(), git.AddFilesAndCommit); err != nil {
		print.Error("Create go mod failed", err)
		return
	}
	bar.Increment()

	// vagrant for docker
	// TODO
	bar.Increment()

	// docker
	// TODO
	bar.Increment()

	// code
	if err := code.Init(kind.Get(), destination.Get(), cfg, golang.Get, git.AddFilesAndCommit); err != nil {
		print.Error("Creating code base failed", err)
		return
	}
	// TODO: service
	bar.Increment()

	// scripts
	if err := scripts.Init(destination.Get(), git.AddFilesAndCommit, git.CreateIgnore); err != nil {
		print.Error("Create scripts failed", err)
		return
	}
	bar.Increment()

	// run goimports to import missing go packages and format code
	if err := golang.GoImports(destination.Get(), git.AddFilesAndCommit); err != nil {
		print.Error("Formatting code failed", err)
		return
	}
	bar.Increment()
	bar.Finish()
}

/*
other TODO:
1. exit with proper Exit Codes
*/
