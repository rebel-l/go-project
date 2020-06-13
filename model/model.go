package model

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/lib/print"
)

type model struct {
	Name       string
	Attributes fields
	rootPath   string
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
		Name:     n,
		rootPath: rootPath,
	}
}

func (m *model) SetID() {
	fmt.Println()
	m.Attributes = append(m.Attributes, NewFieldID())
	fmt.Println()
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

func (m *model) GetSQlTableName() string {
	return strings.ToLower(m.Name) + "s" // TODO: CamelCase to snake_case
}

func (m *model) GetReceiver() string {
	return strings.ToLower(m.Name[0:1])
}

func (m *model) GetStructFields() string {
	return strings.Join(m.getStructFields(0), ", ")
}

func (m *model) GetStructFieldsWithoutID() string {
	return strings.Join(m.getStructFields(1), ", ")
}

func (m *model) getStructFields(start int) []string {
	if len(m.Attributes) < start+1 {
		return nil
	}

	var structFields []string

	for _, v := range m.Attributes[start:] {
		structFields = append(structFields, m.GetReceiver()+"."+v.Name)
	}

	return structFields
}

func (m *model) GetStructFieldsWithIDLast() string {
	structFields := m.getStructFields(1)
	structFields = append(structFields, m.GetReceiver()+"."+m.Attributes[0].Name)

	return strings.Join(structFields, ", ")
}

func (m *model) GetSQLInsert() string {
	numFields := len(m.Attributes) - 1
	fieldNames := m.Attributes.GetSQLFieldNamesWithoutID()

	if m.IsIDUUID() {
		fieldNames = m.Attributes.GetSQLFieldNames()
		numFields++
	}

	placeHolders := make([]string, numFields)
	for i := range placeHolders {
		placeHolders[i] = "?"
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		m.GetSQlTableName(),
		fieldNames,
		strings.Join(placeHolders, ", "),
	)
}

func (m *model) GetExecInsert() string {
	res := "res"
	equal := ":="
	fields := m.GetStructFieldsWithoutID()

	if m.IsIDUUID() {
		res = "_"
		equal = "="
		fields = m.GetStructFields()
	}

	return fmt.Sprintf("%s, err %s db.ExecContext(ctx, q, %s)", res, equal, fields)
}

func (m *model) GetSQLUpdate() string {
	if len(m.Attributes) < 2 {
		return ""
	}

	var fieldNames []string
	for _, v := range m.Attributes[1:] {
		fieldNames = append(fieldNames, v.GetSQlFieldName()+" = ?")
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		m.GetSQlTableName(),
		strings.Join(fieldNames, ", "),
	)
}

func (m *model) GetPackages() ([]string, error) {
	return m.Attributes.GetPackages(m.rootPath)
}

func (m *model) GetIDDefault() string {
	switch m.Attributes[0].FieldType {
	case fieldTypeUUID:
		return "\"\""
	}

	return "0"
}

func (m *model) GetIDEmptyComparison() string {
	return m.Attributes[0].GetEmptyComparison(m.GetReceiver())
}

func (m *model) IsIDUUID() bool {
	return m.Attributes[0].FieldType == fieldTypeUUID
}

func (m *model) GetValidationWithoutID() string {
	return strings.Join(m.Attributes[1:].GetNotNullableFieldsWithComparison(m.GetReceiver()), " || ")
}
