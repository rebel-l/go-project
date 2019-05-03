package destination

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-utils/osutils"
)

var dir string

// Get returns the destination path. If Init() was not executed, it returns empty string.
func Get() string {
	return dir
}

// Init initialises the destination path.
func Init() error {
	if err := detect(); err != nil {
		return setDir()
	}

	printPath()
	c := confirmation()
	if c {
		return nil
	}

	return setDir()
}

func detect() error {
	var err error
	dir, err = os.Getwd()
	return err
}

func printPath() {
	fmt.Printf("The destination path is %s\n", dir)
}

func confirmation() bool {
	t := prompt.Input("Is this path correct? [y/N] ", func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: "y", Description: "confirms directory is correct"},
			{Text: "n", Description: "declines directory and let you choose the right one"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	return strings.ToLower(t) == "y"
}

func setDir() error {
	d := prompt.Input("enter the destination path > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	if osutils.FileOrPathExists(d) {
		dir = d
	} else {
		return fmt.Errorf("path does not exist")
	}
	return nil
}
