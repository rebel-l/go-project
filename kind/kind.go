package kind

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/c-bata/go-prompt"
)

var value string

const (
	kindApplication = "application"
	kindPackage     = "package"
)

// Get returns the kind of the project. If Init() was not called before it returns an empty string.
func Get() string {
	return value
}

// Init request the kind of project from user.
func Init() {
	fmt.Println("Which type of project do you have?")
	for _, v := range possibleKinds() {
		fmt.Printf("%s: %s\n", v.kind, v.description)
	}

	valid := false
	for !valid {
		answer := askUser()
		valid = validate(answer)
		if valid {
			value = answer
		} else {
			errMsg := color.New(color.FgRed, color.Italic)
			_, _ = errMsg.Printf("Project type %s is not valid, please enter again\n", answer)
		}
	}
}

func askUser() string {
	t := prompt.Input("Enter the project type: ", func(d prompt.Document) []prompt.Suggest {
		s := possibleKinds().getSuggestions()
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	return strings.ToLower(t)
}

type kind struct {
	kind        string
	description string
}

func validate(kind string) bool {
	for _, v := range possibleKinds() {
		if kind == v.kind {
			return true
		}
	}

	return false
}
