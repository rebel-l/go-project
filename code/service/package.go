// Package service creates the basic files to start with a new go service
package service

// GetPackages returns a list of packages to install with go mod
func GetPackages() []string {
	// TODO: needs to deal with aliases, e.g. log "github.com/sirupsen/logrus"
	return []string{
		"github.com/rebel-l/go-utils",
	}
}
