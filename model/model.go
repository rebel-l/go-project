package model

import (
	"path"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/rebel-l/go-project/lib/print"
)

type Model struct {
	name            string
	attributes      fields
	destinationPath string
}

func (m Model) createStoreLayer() error {
	return nil
}

func (m Model) createModelLayer() error {
	return nil
}

func (m Model) createMapperLayer() error {
	return nil
}

func (m Model) createCollection() error {
	return nil
}

func NewModel(rootPath string) Model {
	n := prompt.Input("enter the name of your model > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	n = strings.TrimSpace(strings.ToLower(n))
	if n == "" {
		print.Error("model name cannot be empty")
		return NewModel(rootPath)
	}

	return Model{
		name:            n,
		destinationPath: path.Join(rootPath, n),
	}
}
