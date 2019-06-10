package description

import "github.com/c-bata/go-prompt"

var description string

// Get returns the description of the project. If Init() was not called before it returns empty string.
func Get() string {
	return description
}

// Init asks the user to give a project description
func Init() {
	description = prompt.Input("Enter a description to your project: ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	})
}
