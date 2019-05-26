package scripts

import "github.com/rebel-l/go-project/git"

// Init initialises the necessary scripts
func Init(projectPath string, commitCallback git.CallbackAddAndCommit, ignoreCallback git.CallbackCreateIgnore) error {

	return createReportsFolder(projectPath, ignoreCallback)
}
