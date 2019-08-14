// +build !windows

package git

import (
	"fmt"
	"os/exec"
)

func getSetUpstreamCommand(path string) *exec.Cmd {
	command := fmt.Sprintf("cd %s ; git branch --set-upstream-to=origin/master master", path)
	return exec.Command(command)
}
