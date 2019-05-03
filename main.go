package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/destination"
)

func main() {
	fmt.Println()
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("Welcome to Go-Project Tool ...")
	fmt.Println()

	if err := destination.Init(); err != nil {
		printError("Init destination path failed", err)
		return
	}
}

func printError(msg string, err error) {
	errMsg := color.New(color.FgRed, color.Italic)
	_, _ = errMsg.Printf(msg+": %s\n", err)
}
