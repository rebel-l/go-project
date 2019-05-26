package print

import "github.com/fatih/color"

// Warning formats and prints a message as an warning
func Warning(msg string) {
	format := color.New(color.FgHiYellow, color.Bold)
	_, _ = format.Println(msg)
}
