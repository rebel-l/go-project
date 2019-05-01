package workingdir

import (
	"fmt"
	"os"
	"strings"

	"github.com/rebel-l/go-utils/osutils"

	"github.com/c-bata/go-prompt"
)

var dir string

func Get() string {
	return dir
}

func Init() error {
	if err := detect(); err != nil {
		return setDir()
	}

	print()
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

func print() {
	fmt.Printf("The current directory is %s\n", dir)
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
