package model

import (
	"github.com/cheggaaa/pb/v3"
)

type Generator interface {
	Generate(m *model) ([]string, error)
}

type Generators []Generator

func (g Generators) Generate(m *model) ([]string, error) {
	bar := pb.StartNew(len(g))

	var files []string
	for _, v := range g {
		bar.Increment()

		fs, err := v.Generate(m)
		if err != nil {
			return nil, err
		}

		files = append(files, fs...)
	}

	bar.Finish()
	return files, nil
}
