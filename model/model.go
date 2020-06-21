package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/lib/print"
)

const (
	operationCreate = "create"
)

type model struct {
	Name       string
	Attributes fields
	rootPath   string
}

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

func (m *model) GetProjectName() string {
	_, p := filepath.Split(m.rootPath)
	return p
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

func (m *model) GetTestDataCU(operation string) []testDataCRUD {
	// struct nil => covered directly in template
	var testCases []testDataCRUD

	// field only (iterate over all fields)
	for _, f := range m.Attributes {
		if !f.Nullable && (f.Name != fieldNameID || m.Attributes.CountMandatory() > 1) {
			continue // TODO: check for success
		}

		td := testDataCRUD{
			Name:        fmt.Sprintf("%s has %s only", strings.ToLower(m.Name), strings.ToLower(f.Name)),
			Actual:      fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, f.GetTestData()),
			ExpectedErr: fmt.Sprintf("%sstore.ErrDataMissing", strings.ToLower(m.Name)),
		}

		testCases = append(testCases, td)
	}

	// all mandatory fields AND id set ==> CREATE only
	if operation == operationCreate {
		td := testDataCRUD{
			Name:        fmt.Sprintf("%s has id", strings.ToLower(m.Name)),
			Actual:      fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(false, true)),
			ExpectedErr: fmt.Sprintf("%sstore.ErrIDIsSet", strings.ToLower(m.Name)),
		}

		testCases = append(testCases, td)
	}

	// success (a) all fields, (b) only mandatory ==> CREATE (without ID), UPDATE, DELETE, READ
	withID := true
	if operation == operationCreate {
		withID = false
	}
	data := fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(false, withID))
	td := testDataCRUD{
		Name:     fmt.Sprintf("%s has all fields set", strings.ToLower(m.Name)),
		Actual:   data,
		Expected: data,
	}

	testCases = append(testCases, td)

	data = fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(true, withID))
	td = testDataCRUD{
		Name:     fmt.Sprintf("%s has only mandatory fields set", strings.ToLower(m.Name)),
		Actual:   data,
		Expected: data,
	}

	testCases = append(testCases, td)

	// TODO: duplicate (all unique fields seperately) ==> CREATE (without ID), UPDATE
	// TODO: max field length (less, exact, too much) ==> CREATE (without ID), UPDATE
	// TODO: not existing ==> UPDATE

	return testCases
}

func (m *model) GetTestDataRD() []testDataCRUD {
	// struct nil => covered directly in template
	var testCases []testDataCRUD

	// success
	data := fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(false, false))
	td := testDataCRUD{
		Name:     "success",
		Prepare:  data,
		Expected: data,
	}

	testCases = append(testCases, td)

	// not existing
	td = testDataCRUD{
		Name:    "not existing",
		Prepare: fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes[0].GetTestData()),
	}

	testCases = append(testCases, td)

	return testCases
}

func (m *model) GetTestIsValid() []testDataIsValid {
	// struct nil => covered directly in template
	var testCases []testDataIsValid

	// field only (iterate over all fields)
	countMandatory := m.Attributes.CountMandatory()
	for _, f := range m.Attributes {
		expected := "false"
		if !f.Nullable && countMandatory == 1 && f.Name != fieldNameID {
			expected = "true"
		}

		td := testDataIsValid{
			Name:     fmt.Sprintf("%s has %s only", strings.ToLower(m.Name), strings.ToLower(f.Name)),
			Actual:   fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, f.GetTestData()),
			Expected: expected,
		}

		testCases = append(testCases, td)
	}

	// mandatory fields only
	td := testDataIsValid{
		Name:     "mandatory fields only",
		Actual:   fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(true, false)),
		Expected: "true",
	}

	testCases = append(testCases, td)

	td = testDataIsValid{
		Name:     "mandatory fields with id",
		Actual:   fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(true, true)),
		Expected: "true",
	}

	testCases = append(testCases, td)

	// all fields
	td = testDataIsValid{
		Name:     "all fields",
		Actual:   fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(false, true)),
		Expected: "true",
	}

	testCases = append(testCases, td)

	td = testDataIsValid{
		Name:     "all fields without id",
		Actual:   fmt.Sprintf("&%sstore.%s{%s}", strings.ToLower(m.Name), m.Name, m.Attributes.GetTestData(false, false)),
		Expected: "true",
	}

	testCases = append(testCases, td)

	return testCases
}

type testDataCRUD struct {
	Name        string
	Prepare     string
	Actual      string
	Expected    string
	ExpectedErr string
}

type testDataIsValid struct {
	Name     string
	Actual   string
	Expected string
}
