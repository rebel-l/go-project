package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/workingdir"
)

func main() {
	fmt.Println()
	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("Welcome to Go-Project Tool ...")
	fmt.Println()

	if err := workingdir.Init(); err != nil {
		fmt.Printf("Init working directory failed: %s", err)
	}
}
