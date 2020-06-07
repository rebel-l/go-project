package model

import (
	"fmt"
	"path"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/rebel-l/go-project/lib/print"
)

type model struct {
	Name            string
	Attributes      fields
	destinationPath string // TODO: maybe not needed
}

//func (m model) createStoreLayer() error {
//	return nil
//}
//
//func (m model) createModelLayer() error {
//	return nil
//}
//
//func (m model) createMapperLayer() error {
//	return nil
//}
//
//func (m model) createCollection() error {
//	return nil
//}

func NewModel(rootPath string) *model {
	n := prompt.Input("enter the name of your model > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	n = strings.TrimSpace(strings.Title(n))
	if n == "" {
		print.Error("model name cannot be empty")
		return NewModel(rootPath)
	}

	return &model{
		Name:            n,
		destinationPath: path.Join(rootPath, n),
	}
}

func (m *model) AddField() {
	fmt.Println()
	fmt.Println("Add a new field to your model ... leave name empty as you declared all fields")
	f := NewField()
	if f.Name == "" {
		return
	}

	m.Attributes = append(m.Attributes, f)
	m.AddField()
}
