package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/destination"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-project/license"
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
}
