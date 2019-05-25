package kind

import (
	"fmt"

	"github.com/rebel-l/go-utils/option"
)

func possibleKinds() option.Options {
	return option.Options{
		{
			Key:         kindService,
			Description: fmt.Sprintf("creates a project of type %s", kindService),
		},
		{
			Key:         kindPackage,
			Description: fmt.Sprintf("creates a project of type %s", kindPackage),
		},
	}
}
