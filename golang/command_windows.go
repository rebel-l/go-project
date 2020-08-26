// +build windows

package golang

import (
	"os/exec"
)

func getGoModCommand(packageName string) *exec.Cmd {
	return exec.Command("cmd", "/C", "go", "mod", "init", packageName)
}

func getGoGetCommand(packageName string) *exec.Cmd {
	return exec.Command("cmd", "/C", "go", "get", packageName)
}

func getGoModTidyCommand() *exec.Cmd {
	return exec.Command("cmd", "/C", "go", "mod", "tidy")
}

func getGoImportsCommand(path string) *exec.Cmd {
	return exec.Command("cmd", "/C", "goimports", "-w", path)
}
