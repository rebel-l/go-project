package options

import (
	"fmt"

	"github.com/rebel-l/go-utils/option"
)

// Print prints the options
func Print(options option.Options) {
	_ = options.ForAll(func(option option.Option) error {
		fmt.Printf("%s: %s\n", option.Key, option.Description)
		return nil
	})
}
