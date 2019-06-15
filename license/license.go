package license

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/rebel-l/go-project/git"

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
var prefix string

// Get retuns the name selected license. If Init() was not called before it returns an empty string.
func Get() string {
	return value
}

// GetPrefix returns the license prefix.
func GetPrefix() string {
	return prefix
}

// Init let the user select the license and creates license file
func Init(path string, author string, projectDescription string, commit git.CallbackAddAndCommit) error {
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

	err := createLicense(filename, newParameters(author, projectDescription))
	if err != nil {
		return err
	}

	return commit([]string{filename}, "added license")
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
	t := prompt.Input("Enter the license: ", func(d prompt.Document) []prompt.Suggest {
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

	if err = tmpl.ExecuteTemplate(file, value, params); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.ExecuteTemplate(&buf, fmt.Sprintf("%s_prefix", value), params); err == nil { // TODO: do also error handling in a proper way. Ignore only error that template was not found, as this is expected
		prefix = buf.String()
	}
	return nil
}

type parameters struct {
	Year               int
	Author             string
	ProjectDescription string
}

func newParameters(author, projectDescription string) parameters {
	return parameters{
		Year:               time.Now().Year(),
		Author:             author,
		ProjectDescription: projectDescription,
	}
}
