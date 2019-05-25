package license

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/lib/options"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-utils/option"
	"github.com/rebel-l/go-utils/osutils"
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
func Init(path string, author string, commit func(files []string, msg string) error) error {
	filename := filepath.Join(path, "LICENSE")
	if osutils.FileOrPathExists(filename) {
		print.Info("Skip creating a license file as it already exists")
		return nil
	}

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
			print.Error(fmt.Sprintf("License %s is not valid, please enter again\n", answer))
		}
	}

	if value == strings.ToLower(licenseNone) {
		return nil
	}

	err := createLicense(filename, newParameters(author))
	if err != nil {
		return err
	}

	return commit([]string{"LICENSE"}, "added license")
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

func createLicense(filename string, params parameters) error {
	pattern := filepath.Join("./license/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create license file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	return tmpl.ExecuteTemplate(file, value, params)
}

type parameters struct {
	Year   int
	Author string
}

func newParameters(author string) parameters {
	return parameters{
		Year:   time.Now().Year(),
		Author: author,
	}
}

// TODO: get license prefix for source code: GPLv3
