package model

func Init(path string) error {
	m := NewModel(path)

	g := getGenerators()

	return g.Generate(m)
}

func getGenerators() Generators {
	var g Generators

	g = append(g, &SQL{})

	return g
}
