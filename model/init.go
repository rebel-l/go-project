package model

import (
	"fmt"
	"strings"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/golang"
)

const (
	packageGoUtils = "github.com/rebel-l/go-utils@v1.2.0-rc.3" // TODO: remove version number
)

func Init(path string, commit git.CallbackAddAndCommit) error {
	if err := golang.Get(path, packageGoUtils); err != nil {
		return err
	}

	m := NewModel(path)
	m.SetID()
	m.AddField()

	g := getGenerators(path)

	fmt.Println()

	files, err := g.Generate(m)
	if err != nil {
		return err
	}

	// clean paths from files as they are added from commit function
	for i, v := range files {
		files[i] = strings.Replace(v, path, "", 1)
	}

	if err := commit(files, fmt.Sprintf("added model %s", m.Name), 1); err != nil {
		return err
	}

	return golang.GoImports(path, commit, 2)
}

func getGenerators(path string) Generators {
	var g Generators

	g = append(g, &config{rootPath: path})
	g = append(g, &bootstrap{rootPath: path})
	g = append(g, &sql{rootPath: path})
	g = append(g, &store{rootPath: path})
	g = append(g, &modelGen{rootPath: path})
	g = append(g, &mapper{rootPath: path})

	return g
}
