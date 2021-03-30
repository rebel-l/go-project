package git

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func GetPackage(path string) (string, error) {
	var err error
	repo, err = git.PlainOpen(path) // it's important that the global variable "repo" is used
	if err != nil {
		return "", err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", err
	}

	for _, r := range remotes {
		if r.Config() != nil {
			for _, u := range r.Config().URLs {
				if len(u) > 0 {
					return getReplacers().get(u), nil
				}
			}
		}
	}

	return "", nil
}

type replacer struct {
	old string
	new string
}

type replacers []replacer

func (r replacers) get(s string) string {
	for _, v := range r {
		s = strings.Replace(s, v.old, v.new, -1)
	}

	return s
}

func getReplacers() replacers {
	return replacers{
		{
			old: ".git",
			new: "",
		},
		{
			old: "https://",
			new: "",
		},
		{
			old: "git@",
			new: "",
		},
		{
			old: ":",
			new: "/",
		},
	}
}
