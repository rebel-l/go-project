package model

import "fmt"

func Init(path string) error {
	m := NewModel(path)
	m.AddField()

	fmt.Printf("%#v\n", m)
	fmt.Printf("%#v\n", m.attributes)

	g := getGenerators(path)

	fmt.Println()

	return g.Generate(m)
}

func getGenerators(path string) Generators {
	var g Generators

	g = append(g, &sql{rootPath: path})

	return g
}
