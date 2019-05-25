package options

import (
	"github.com/c-bata/go-prompt"
	"github.com/rebel-l/go-utils/option"
)

// GetSuggestions returns suggestions based on the given options
func GetSuggestions(options option.Options) []prompt.Suggest {
	s := make([]prompt.Suggest, len(options))
	for i, v := range options {
		s[i] = prompt.Suggest{
			Text:        v.Key,
			Description: v.Description,
		}
	}
	return s
}
