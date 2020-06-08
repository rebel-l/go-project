package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rebel-l/go-project/lib/options"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-utils/option"

	"github.com/c-bata/go-prompt"
)

const (
	fieldTypeString = "string"
	fieldTypeInt    = "int"
	fieldTypeFloat  = "float"
	fieldTypeTime   = "time"
	fieldTypeBool   = "bool"
	fieldTypeUUID   = "uuid"
)

type fields []*field

func (f fields) IsLast(index int) bool {
	return index >= len(f)-1
}

type field struct {
	Name         string
	PrimaryKey   bool
	FieldType    string
	DefaultValue string
	MaxLength    int
	Nullable     bool
	Unique       bool
}

func NewField() *field {
	f := &field{}

	f.setName()
	if f.Name == "" {
		return f
	}

	f.setFieldType()
	f.setMaxLength()
	f.setUnique()
	f.setDefaultValue()
	f.setNullable()

	return f
}

func (f *field) setName() {
	n := prompt.Input("enter the name of the field > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	n = strings.TrimSpace(n)
	if strings.ToLower(n) == "id" {
		f.PrimaryKey = true
		f.Name = strings.ToUpper(n)
	} else {
		f.Name = strings.Title(n)
	}
}

func (f *field) setFieldType() {
	fieldTypes := getPossibleFieldTypes()
	t := prompt.Input("enter the the type of the field > ", func(d prompt.Document) []prompt.Suggest {
		s := options.GetSuggestions(fieldTypes)
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}, prompt.OptionInputTextColor(prompt.Yellow))

	t = strings.TrimSpace(strings.ToLower(t))

	if !fieldTypes.IsValidOptionCI(t) {
		var keys []string
		_ = fieldTypes.ForAll(func(option option.Option) error {
			keys = append(keys, option.Key)
			return nil
		})

		print.Error(
			fmt.Sprintf("field type must be one of the following values: %s", strings.Join(keys, ", ")),
		)
		f.setFieldType()
	}

	f.FieldType = t
}

func (f *field) setDefaultValue() {
	if f.PrimaryKey || f.Unique {
		return
	}

	f.DefaultValue = prompt.Input("enter the default value of the field > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))
}

func (f *field) setNullable() {
	if f.PrimaryKey || f.Unique {
		return
	}

	n := prompt.Input("is the field nullable? [y/N] > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	if strings.TrimSpace(strings.ToLower(n)) == "y" {
		f.Nullable = true
	}
}

func (f *field) setUnique() {
	if f.PrimaryKey {
		return
	}

	n := prompt.Input("should the value of the field be unique? [y/N] > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	if strings.TrimSpace(strings.ToLower(n)) == "y" {
		f.Unique = true
	}
}

func (f *field) setMaxLength() {
	if f.FieldType != fieldTypeString {
		return
	}

	l := prompt.Input("enter the maximum length of the fields value > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	var err error
	f.MaxLength, err = strconv.Atoi(l)
	if err != nil {
		print.Error("maximum length must be a valid integer", err)
		f.setMaxLength()
	}
}

func getPossibleFieldTypes() option.Options {
	return option.Options{
		{
			Key:         fieldTypeString,
			Description: "value of type string",
		},
		{
			Key:         fieldTypeInt,
			Description: "value of type integer",
		},
		{
			Key:         fieldTypeFloat,
			Description: "value of type float",
		},
		{
			Key:         fieldTypeTime,
			Description: "value of type time",
		},
		{
			Key:         fieldTypeBool,
			Description: "value of type bool",
		},
		{
			Key:         fieldTypeUUID,
			Description: "value of type uuid",
		},
	}
}
