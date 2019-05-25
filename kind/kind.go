package kind

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/lib/options"
	"github.com/rebel-l/go-project/lib/print"
)

var value string

const (
	kindService = "service"
	kindPackage = "package"
)

// Get returns the kind of the project. If Init() was not called before it returns an empty string.
func Get() string {
	return value
}

// Init request the kind of project from user.
func Init() {
	fmt.Println("Which type of project do you have?")
	kinds := possibleKinds()
	options.Print(kinds)

	valid := false
	for !valid {
		answer := askUser()
		valid = kinds.IsValidOption(answer)
		if valid {
			value = answer
		} else {
			print.Error(fmt.Sprintf("Project type %s is not valid, please enter again\n", answer))
		}
	}
}

func askUser() string {
	t := prompt.Input("Enter the project type: ", func(d prompt.Document) []prompt.Suggest {
		s := options.GetSuggestions(possibleKinds())
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	return strings.ToLower(t)
}
