package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"

	"github.com/google/uuid"

	"github.com/rebel-l/go-project/lib/options"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-utils/option"
	"github.com/rebel-l/go-utils/randutils"

	"github.com/c-bata/go-prompt"
)

const (
	fieldNameID = "ID"

	fieldTypeString    = "string"
	fieldTypeEmail     = "email"
	fieldTypeFirstName = "firstname"
	fieldTypeLastName  = "lastname"
	fieldTypeInt       = "int"
	fieldTypeFloat     = "float"
	fieldTypeTime      = "time"
	fieldTypeBool      = "bool"
	fieldTypeUUID      = "uuid"

	packageUUID = "github.com/google/uuid"
)

func (f *field) GetGoFieldType() string {
	var ft string
	switch f.FieldType {
	case fieldTypeUUID:
		ft = "uuid.UUID"
	case fieldTypeFirstName,
		fieldTypeLastName,
		fieldTypeEmail:
		ft = "string"
	default:
		ft = f.FieldType
	}

	return ft
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

func (f *field) GetStoreField() string {
	return fmt.Sprintf("%s	%s	`db:\"%s\"`", f.Name, f.GetGoFieldType(), f.GetSQlFieldName())
}

func (f *field) GetSQlFieldName() string {
	return strings.ToLower(f.Name) // TODO: CamelCase to snake_case
}

func (f *field) GetSQLField() string { // TODO: support sql dialect ... maybe with parameter
	line := strings.ToLower(f.Name)

	switch f.FieldType {
	case fieldTypeUUID:
		line += " CHAR(36)"
	case fieldTypeString,
		fieldTypeEmail,
		fieldTypeFirstName,
		fieldTypeLastName:
		line += fmt.Sprintf(" VARCHAR(%d)", f.MaxLength)
	case fieldTypeBool,
		fieldTypeInt:
		line += " INTEGER"
	case fieldTypeFloat:
		line += " REAL"
	case fieldTypeTime:
		line += " DATETIME"
	}

	if !f.Nullable {
		line += " NOT NULL"
	}

	if f.DefaultValue != "" {
		switch f.FieldType {
		case fieldTypeString,
			fieldTypeEmail,
			fieldTypeFirstName,
			fieldTypeLastName:
			line += fmt.Sprintf(" DEFAULT '%s'", f.DefaultValue)
		case fieldTypeBool:
			if strings.ToLower(f.DefaultValue) == "true" {
				line += " DEFAULT 1"
			} else {
				line += " DEFAULT 0"
			}
		default:
			line += fmt.Sprintf(" DEFAULT %s", f.DefaultValue)
		}
	}

	if f.PrimaryKey {
		line += " PRIMARY KEY"
		if f.FieldType == fieldTypeInt {
			line += " AUTOINCREMENT"
		}
	}

	return line
}

func (f *field) GetEmptyComparison(receiver string) string {
	operator := "=="

	field := receiver + "." + f.Name
	switch f.FieldType {
	case fieldTypeUUID:
		return fmt.Sprintf("uuidutils.IsEmpty(%s)", field)
	case fieldTypeTime:
		field += ".IsZero()"
		operator = ""
	}

	return fmt.Sprintf("%s %s %s", field, operator, f.GetDefaultValue())
}

func (f *field) GetDefaultValue() string {
	var value string
	switch f.FieldType {
	case fieldTypeUUID,
		fieldTypeString,
		fieldTypeEmail,
		fieldTypeFirstName,
		fieldTypeLastName:
		value = "\"\""
	case fieldTypeBool,
		fieldTypeInt:
		value = "0"
	case fieldTypeFloat:
		value = "0.0"
	}

	return value
}

func (f *field) GetFormat() string {
	var value string
	switch f.FieldType {
	case fieldTypeUUID,
		fieldTypeString,
		fieldTypeEmail,
		fieldTypeFirstName,
		fieldTypeLastName:
		value = "%s"
	case fieldTypeBool:
		value = "%t"
	case fieldTypeInt:
		value = "%d"
	case fieldTypeFloat:
		value = "%f"
	}

	return value
}

func (f *field) GetTestData() string {
	data := f.Name + ": "

	switch f.FieldType {
	case fieldTypeUUID:
		u, err := uuid.NewRandom()
		if err != nil {
			return ""
		}
		data += fmt.Sprintf("testingutils.UUIDParse(t, \"%s\")", u.String())
	case fieldTypeString:
		max := 50
		if f.MaxLength > 0 {
			max = f.MaxLength
		}
		data += fmt.Sprintf("\"%s\"", randomdata.RandStringRunes(randutils.Int(5, max)))
	case fieldTypeEmail:
		data += fmt.Sprintf("\"%s\"", randomdata.Email())
	case fieldTypeFirstName:
		data += fmt.Sprintf("\"%s\"", randomdata.FirstName(randomdata.RandomGender))
	case fieldTypeLastName:
		data += fmt.Sprintf("\"%s\"", randomdata.LastName())
	case fieldTypeInt:
		data += fmt.Sprintf("%d", randutils.Int(1, math.MaxInt16))
	case fieldTypeFloat:
		data += fmt.Sprintf("%f", randomdata.Decimal(10, 10000))
	case fieldTypeTime:
		data += fmt.Sprintf("time.Parse(\"\\\"2006-01-02 15:04:05.999999999 -0700 MST\\\"\", %s)", time.Now().String())
	case fieldTypeBool:
		data += " true"
	}

	return data
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

func NewFieldID() *field {
	f := &field{
		Name:       fieldNameID,
		PrimaryKey: true,
	}

	fieldTypes := getPossibleIDTypes()
	t := prompt.Input("enter the type of the ID > ", func(d prompt.Document) []prompt.Suggest {
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
			fmt.Sprintf("type of ID must be one of the following values: %s", strings.Join(keys, ", ")),
		)
		f.setFieldType()
	}

	f.FieldType = t

	return f
}

func (f *field) setName() {
	n := prompt.Input("enter the name of the field > ", func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	}, prompt.OptionInputTextColor(prompt.Yellow))

	n = strings.TrimSpace(n)
	f.Name = strings.Title(n)
}

func (f *field) setFieldType() {
	fieldTypes := getPossibleFieldTypes()
	t := prompt.Input("enter the type of the field > ", func(d prompt.Document) []prompt.Suggest {
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
	if f.FieldType != fieldTypeString &&
		f.FieldType != fieldTypeEmail &&
		f.FieldType != fieldTypeFirstName &&
		f.FieldType != fieldTypeLastName {
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
			Key:         fieldTypeEmail,
			Description: "value of type email which is at the end just a string",
		},
		{
			Key:         fieldTypeFirstName,
			Description: "value of type first name which is at the end just a string",
		},
		{
			Key:         fieldTypeLastName,
			Description: "value of type last name which is at the end just a string",
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

func getPossibleIDTypes() option.Options {
	return option.Options{
		{
			Key:         fieldTypeInt,
			Description: "value of type integer",
		},
		{
			Key:         fieldTypeUUID,
			Description: "value of type uuid",
		},
	}
}
