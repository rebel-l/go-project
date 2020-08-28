// +build !windows

package golang

import "os/exec"

func getGoModCommand(packageName string) *exec.Cmd {
	return exec.Command("go", "mod", "init", packageName)
}

func getGoGetCommand(packageName string) *exec.Cmd {
	return exec.Command("go", "get", packageName)
}

func getGoModTidyCommand() *exec.Cmd {
	return exec.Command("go", "mod", "tidy")
}
