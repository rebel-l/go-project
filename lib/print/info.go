package print

import "github.com/fatih/color"

// Info formats and prints a message as an information
func Info(msg string) {
	format := color.New(color.FgHiGreen)
	_, _ = format.Println(msg)
}
