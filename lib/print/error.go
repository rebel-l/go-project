package print

import (
	"fmt"

	"github.com/fatih/color"
)

var errorFmt = []color.Attribute{color.FgRed, color.Bold}

// Error formats and prints a message as an error
func Error(msg string, err ...error) {
	format := color.New(errorFmt...)
	if len(err) == 0 {
		_, _ = format.Println(msg)
		return
	}

	var errMsg string
	for _, e := range err {
		errMsg += fmt.Sprintf("%s\n", e.Error())
	}
	_, _ = format.Printf("%s: %s\n", msg, errMsg)
}
