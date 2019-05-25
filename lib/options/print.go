package options

import (
	"github.com/fatih/color"

	"github.com/rebel-l/go-utils/option"
)

// Print prints the options
func Print(options option.Options) {
	_ = options.ForAll(func(option option.Option) error {
		msg := color.New(color.FgCyan, color.Italic)
		_, _ = msg.Printf("%s: %s\n", option.Key, option.Description)
		return nil
	})
}
