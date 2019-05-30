package readme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rebel-l/go-project/git"
)

// Init intialises th readme file
func Init(projectPath string, repository string, license string, commit git.CallbackAddAndCommit) error {
	pattern := filepath.Join("./readme/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filename := filepath.Join(projectPath, "README.md")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create readme file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, "readme", extractParams(repository, license)); err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}
	return commit([]string{filename}, "added readme")
}

type parameters struct {
	Project     string
	GitDomain   string
	GitUsername string
	License     string
}

func (p parameters) GetGitCompany() string {
	return strings.Split(p.GitDomain, ".")[0]
}

func extractParams(repository string, license string) parameters {
	/*
		Example strings to split:
			https://github.com/rebel-l/auth-service.git
			git@github.com:rebel-l/auth-service.git
	*/
	params := parameters{License: license}
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
