package main

import (
	"fmt"

	"github.com/rebel-l/go-project/license"

	"github.com/rebel-l/go-project/git"

	"github.com/rebel-l/go-project/kind"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/destination"
)

func main() {
	// introduction
	fmt.Println()
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("Welcome to Go-Project Tool ...")
	fmt.Println()

	// destination path
	if err := destination.Init(); err != nil {
		printError("Init destination path failed", err)
		return
	}
	fmt.Println()

	// project kind
	kind.Init()
	fmt.Println()

	// git setup
	git.Setup(destination.Get())

	// license
	if err := license.Init(destination.Get()); err != nil {
		printError("Init license failed", err)
		return
	}
	fmt.Println()
}

func printError(msg string, err error) {
	errMsg := color.New(color.FgRed, color.Italic)
	_, _ = errMsg.Printf(msg+": %s\n\n", err)
}
