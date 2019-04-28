package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"

	"github.com/rebel-l/go-project/workingdir"
)

func main() {
	fmt.Println()

	title := color.New(color.Bold, color.FgGreen)
	_, _ = title.Println("Welcome to Go-Project Tool ...")

	if err := workingdir.Init(); err != nil {
		fmt.Printf("Init working directory failed: %s", err)
	}

	fmt.Printf("The current directory is %s", workingdir.Get())

	fmt.Println("Please select table.")
	t := prompt.Input("> ", completer)
	fmt.Println("You selected " + t)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
