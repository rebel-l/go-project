package config

import (
	"fmt"
	"strings"
)

// Config represents configuration data
type Config struct {
	Project       string
	GitDomain     string
	GitUsername   string
	License       string
	LicensePrefix string
	Description   string
}

// GetGitCompany returns the GitDomain without top level domain
func (c Config) GetGitCompany() string {
	return strings.Split(c.GitDomain, ".")[0]
}

// GetPackage returns the full package name combined by GitDomain/GitUsername/Project
func (c Config) GetPackage() string {
	return fmt.Sprintf("%s/%s/%s", c.GitDomain, c.GitUsername, c.Project)
}

// New returns a new config extracted the repository and license
func New(repository, license, licensePrefix, description string) Config {
	/*
		Example strings to split:
			https://github.com/rebel-l/auth-service.git
			git@github.com:rebel-l/auth-service.git
	*/
	params := Config{License: license, LicensePrefix: licensePrefix, Description: description}
	repository = strings.ToLower(repository)
	pieces := strings.Split(repository, "/")
	params.Project = strings.Replace(pieces[len(pieces)-1], ".git", "", -1)

	switch len(pieces) {
	case 2:
		sub := strings.Split(pieces[0], ":")
		if len(sub) == 2 {
			params.GitDomain = strings.Replace(sub[0], "git@", "", -1)
			params.GitUsername = sub[1]
		}
	case 5:
		params.GitUsername = pieces[3]
		params.GitDomain = pieces[2]
	}

	return params
}
