package model

import (
	"fmt"

	"github.com/rebel-l/go-project/golang"
)

const (
	packageGoUtils = "github.com/rebel-l/go-utils"
)

func Init(path string) error {
	if err := golang.Get(path, packageGoUtils); err != nil {
		return err
	}

	m := NewModel(path)
	m.SetID()
	m.AddField()

	g := getGenerators(path)

	fmt.Println()

	return g.Generate(m)
}

func getGenerators(path string) Generators {
	var g Generators

	g = append(g, &config{rootPath: path})
	g = append(g, &bootstrap{rootPath: path})
	g = append(g, &sql{rootPath: path})
	g = append(g, &store{rootPath: path})

	return g
}
