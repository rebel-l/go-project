package git

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func GetPackage(path string) (string, error) {
	repo, err := git.PlainOpen(path)
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
					return strings.Replace(
						strings.Replace(u, ".git", "", -1),
						"https://",
						"",
						-1), nil
				}
			}
		}
	}

	return "", nil
}
