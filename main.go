package main

import (
	"fmt"

	"github.com/rebel-l/go-project/workingdir"
)

func main() {
	fmt.Println()
	fmt.Println("Welcome to Go-Project Tool ...")

	if err := workingdir.Init(); err != nil {
		fmt.Printf("Init working directory failed: %s", err)
	}

	fmt.Printf("The current directory is %s", workingdir.Get())
}
