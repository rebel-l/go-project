package kind

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

type kinds []kind

func (k kinds) getSuggestions() []prompt.Suggest {
	s := make([]prompt.Suggest, len(k))
	for i, v := range k {
		s[i] = prompt.Suggest{
			Text:        v.kind,
			Description: v.description,
		}
	}
	return s
}

func possibleKinds() kinds {
	return kinds{
		{
			kind:        kindApplication,
			description: fmt.Sprintf("creates a project of type %s", kindApplication),
		},
		{
			kind:        kindPackage,
			description: fmt.Sprintf("creates a project of type %s", kindPackage),
		},
	}
}
