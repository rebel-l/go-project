package model

import "fmt"

func Init(path string) error {
	m := NewModel(path)
	m.AddField()

	g := getGenerators(path)

	fmt.Println()

	return g.Generate(m)
}

func getGenerators(path string) Generators {
	var g Generators

	g = append(g, &sql{rootPath: path})
	g = append(g, &store{rootPath: path})

	return g
}