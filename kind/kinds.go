package kind

import (
	"fmt"

	"github.com/rebel-l/go-utils/option"
)

func possibleKinds() option.Options {
	return option.Options{
		{
			Key:         Service,
			Description: fmt.Sprintf("creates a project of type %s", Service),
		},
		{
			Key:         Package,
			Description: fmt.Sprintf("creates a project of type %s", Package),
		},
	}
}
