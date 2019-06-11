// Package service creates the basic files to start with a new go service
package service

import (
	"github.com/rebel-l/go-project/golang"
)

// GetPackages returns a list of packages to install with go mod
func GetPackages() golang.Imports {
	return golang.Imports{
		{Name: "github.com/rebel-l/go-utils"},
		{Name: "github.com/sirupsen/logrus", Alias: "log"},
	}
}
