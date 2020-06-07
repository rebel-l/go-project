package model

import (
	"github.com/cheggaaa/pb/v3"
)

type Generator interface {
	Generate(m *model) error
}

type Generators []Generator

func (g Generators) Generate(m *model) error {
	bar := pb.StartNew(len(g))
	for _, v := range g {
		bar.Increment()
		if err := v.Generate(m); err != nil {
			return err
		}
	}

	bar.Finish()
	return nil
}