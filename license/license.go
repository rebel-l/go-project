package license

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/fatih/color"

	"github.com/rebel-l/go-project/lib/options"
	"github.com/rebel-l/go-utils/option"
)

const (
	licenseGPL3 = "GPLv3"
	licenseMIT  = "MIT"
	licenseNone = "none"
)

var value string
var description = "Creates project under the %s license"

// Get retuns the name selected license. If Init() was not called before it returns an empty string.
func Get() string {
	return value
}

// Init let the user select the license and creates license file
func Init() {
	fmt.Println("Under which license should this project be published?")
	licenses := getPossibleLicenses()
	options.Print(licenses)

	valid := false
	for !valid {
		answer := askUser()
		valid = licenses.IsValidOptionCI(answer)
		if valid {
			value = answer
		} else {
			errMsg := color.New(color.FgRed, color.Italic)
			_, _ = errMsg.Printf("License %s is not valid, please enter again\n", answer)
		}
	}
}

func getPossibleLicenses() option.Options {
	return option.Options{
		{
			Key:         licenseGPL3,
			Description: fmt.Sprintf(description, licenseGPL3),
		},
		{
			Key:         licenseMIT,
			Description: fmt.Sprintf(description, licenseMIT),
		},
		{
			Key:         licenseNone,
			Description: "skips the creation of any license, if you want to do it later you need to do it manually",
		},
	}
}

func askUser() string {
	t := prompt.Input("Enter the project type: ", func(d prompt.Document) []prompt.Suggest {
		s := options.GetSuggestions(getPossibleLicenses())
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	return strings.ToLower(t)
}
