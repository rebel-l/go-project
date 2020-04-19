package dialog

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

func Confirmation(question string) bool {
	answer := "y"
	t := prompt.Input(fmt.Sprintf("%s [Y/n] ", question), func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if t != "" {
		answer = t
	}

	return strings.ToLower(answer) == "y"
}
