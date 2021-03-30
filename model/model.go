package model

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/print"
)

const (
	operationCreate = "create"
	operationUpdate = "update"
	operationRead   = "read"

	packageBootstrap      = "bootstrap"
	packageBootstrapTest  = "bootstrap_test"
	packageConfigTest     = "config_test"
	packageTypeStore      = "store"
	packageTypeStoreTest  = "store_test"
	packageTypeModel      = "model"
	packageTypeModelTest  = "model_test"
	packageTypeMapper     = "mapper"
	packageTypeMapperTest = "mapper_test"
)

type model struct {
	Name       string
	Attributes fields
	rootPath   string
}

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

func (m *model) GetImports(packageType string) []string {
	var packages []string

	name := strings.ToLower(m.Name)

	p, err := git.GetPackage(m.rootPath)
	if err != nil {
		return packages
	}

	importPrefix := fmt.Sprintf("%s/%s/%s", p, name, name)

	switch strings.ToLower(packageType) {
	case packageBootstrap:
		packages = append(
			packages,
			"github.com/jmoiron/sqlx",
			p+"/config",
			"github.com/rebel-l/go-utils/osutils",
			"github.com/rebel-l/schema",
		)
	case packageBootstrapTest:
		packages = append(
			packages,
			p+"/bootstrap",
			p+"/config",
			"github.com/rebel-l/go-utils/osutils",
			"github.com/rebel-l/go-utils/slice",
		)
	case packageConfigTest:
		packages = append(
			packages,
			p+"/config",
		)
	case packageTypeMapper:
		if m.GetIDType() == "uuid.UUID" {
			packages = append(packages, "github.com/google/uuid", "github.com/rebel-l/go-utils/uuidutils")
		}

		packages = append(
			packages,
			"github.com/jmoiron/sqlx",
			importPrefix+"model",
			importPrefix+"store",
		)
	case packageTypeMapperTest:
		if m.GetIDType() == "uuid.UUID" {
			packages = append(packages, "github.com/google/uuid")
		}

		packages = append(
			packages,
			"github.com/jmoiron/sqlx",
			"github.com/rebel-l/go-utils/osutils",
			"github.com/rebel-l/go-utils/testingutils",
			importPrefix+"mapper",
			importPrefix+"model",
			p+"/bootstrap",
			p+"/config",
		)
	case packageTypeModel:
		if m.GetIDType() == "uuid.UUID" {
			packages = append(packages, "github.com/google/uuid")
		}
	case packageTypeModelTest:
		packages = append(packages, "github.com/rebel-l/go-utils/testingutils", importPrefix+"model")
	case packageTypeStore:
		if m.GetIDType() == "uuid.UUID" {
			packages = append(packages, "github.com/google/uuid", "github.com/rebel-l/go-utils/uuidutils")
		}

		packages = append(packages, "github.com/jmoiron/sqlx")
	case packageTypeStoreTest:
		if m.GetIDType() == "uuid.UUID" {
			packages = append(packages, "github.com/google/uuid", "github.com/rebel-l/go-utils/uuidutils")
		}

		packages = append(
			packages,
			"github.com/jmoiron/sqlx",
			"github.com/rebel-l/go-utils/osutils",
			"github.com/rebel-l/go-utils/testingutils",
			importPrefix+"store",
			p+"/bootstrap",
			p+"/config",
		)
	}

	sort.Strings(packages)

	return packages
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
		`
		INSERT INTO %s (%s) 
		VALUES (%s);
	`,
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
		`
		UPDATE %s 
		SET %s 
		WHERE id = ?;
	`,
		m.GetSQlTableName(),
		strings.Join(fieldNames, ", "),
	)
}

func (m *model) GetPackages() ([]string, error) {
	// TODO: maybe this is deprecated
	return m.Attributes.GetPackages(m.rootPath)
}

func (m *model) GetIDDefault() string {
	switch m.Attributes[0].FieldType {
	case fieldTypeUUID:
		return "\"\""
	}

	return "0"
}

func (m *model) GetIDType() string {
	return m.Attributes[0].GetGoFieldType()
}

func (m *model) GetIDEmptyComparison(receiver string) string {
	if receiver == "" {
		receiver = m.GetReceiver()
	}

	return m.Attributes[0].GetEmptyComparison(receiver)
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
	for i, f := range m.Attributes {
		if !f.Nullable && (f.Name != fieldNameID || m.Attributes.CountMandatory() > 1) {
			continue // TODO: check for success
		}

		nm := m.Clone()
		nm.Attributes.ResetTestData()
		nm.Attributes[i].GetTestData()
		td := testDataCRUD{
			Name:        fmt.Sprintf("%s has %s only", strings.ToLower(m.Name), strings.ToLower(f.Name)),
			Actual:      nm,
			ExpectedErr: fmt.Sprintf("%sstore.ErrDataMissing", strings.ToLower(m.Name)),
		}

		testCases = append(testCases, td)
	}

	// all mandatory fields AND id set ==> CREATE only
	if operation == operationCreate {
		nm := m.GetTestDataForStruct(false, true)
		td := testDataCRUD{
			Name:        fmt.Sprintf("%s has id", strings.ToLower(m.Name)),
			Actual:      nm,
			ExpectedErr: fmt.Sprintf("%sstore.ErrIDIsSet", strings.ToLower(m.Name)),
		}

		testCases = append(testCases, td)
	}

	// all mandatory fields set AND id is missing ==> UPDATE only
	if operation == operationUpdate {
		nm := m.GetTestDataForStruct(false, false)
		td := testDataCRUD{
			Name:        fmt.Sprintf("%s has no id", strings.ToLower(m.Name)),
			Actual:      nm,
			ExpectedErr: fmt.Sprintf("%sstore.ErrIDMissing", strings.ToLower(m.Name)),
		}

		testCases = append(testCases, td)

		// not existing ==> UPDATE
		nm = m.GetTestDataForStruct(true, true)
		td = testDataCRUD{
			Name:        "not existing",
			Actual:      nm,
			ExpectedErr: "sql.ErrNoRows",
		}

		testCases = append(testCases, td)
	}

	// success (a) all fields, (b) only mandatory ==> CREATE (without ID), UPDATE, DELETE, READ
	nm := m.GetTestDataForStruct(false, false)
	td := testDataCRUD{
		Name:     fmt.Sprintf("%s has all fields set", strings.ToLower(m.Name)),
		Actual:   nm,
		Expected: nm,
	}

	if operation == operationUpdate {
		td.Prepare = m.GetTestDataForStruct(false, false)
	}

	testCases = append(testCases, td)

	nm = m.GetTestDataForStruct(true, false)
	td = testDataCRUD{
		Name:     fmt.Sprintf("%s has only mandatory fields set", strings.ToLower(m.Name)),
		Actual:   nm,
		Expected: nm,
	}

	if operation == operationUpdate {
		td.Prepare = m.GetTestDataForStruct(true, false)
	}

	testCases = append(testCases, td)

	// TODO: duplicate (all unique fields seperately) ==> CREATE (without ID), UPDATE
	// TODO: max field length (less, exact, too much) ==> CREATE (without ID), UPDATE

	return testCases
}

func (m *model) GetTestDataRD(operation string) []testDataCRUD {
	// struct nil => covered directly in template
	var testCases []testDataCRUD

	// success
	nm := m.GetTestDataForStruct(false, false)
	td := testDataCRUD{
		Name:     "success",
		Prepare:  nm,
		Expected: nm,
	}

	testCases = append(testCases, td)

	// not existing
	nm = m.Clone()
	nm.Attributes.ResetTestData()
	_ = nm.Attributes[0].GetTestData()
	td = testDataCRUD{
		Name:    "not existing",
		Prepare: nm,
	}

	if operation == operationRead {
		td.ExpectedErr = "sql.ErrNoRows"
	}

	testCases = append(testCases, td)

	return testCases
}

func (m *model) GetTestIsValid() []testDataIsValid {
	// struct nil => covered directly in template
	var testCases []testDataIsValid

	// field only (iterate over all fields)
	countMandatory := m.Attributes.CountMandatory()
	for i, f := range m.Attributes {
		expected := "false"
		if !f.Nullable && countMandatory == 1 && f.Name != fieldNameID {
			expected = "true"
		}

		nm := m.Clone()
		nm.Attributes.ResetTestData()
		nm.Attributes[i].GetTestData()
		td := testDataIsValid{
			Name:     fmt.Sprintf("%s has %s only", strings.ToLower(m.Name), strings.ToLower(f.Name)),
			Actual:   nm,
			Expected: expected,
		}

		testCases = append(testCases, td)
	}

	// mandatory fields only
	nm := m.Clone()
	nm.Attributes.ResetTestData()
	nm.Attributes.GetTestDataForStruct(true, false)
	td := testDataIsValid{
		Name:     "mandatory fields only",
		Actual:   nm,
		Expected: "true",
	}

	testCases = append(testCases, td)

	nm = m.Clone()
	nm.Attributes.ResetTestData()
	nm.Attributes.GetTestDataForStruct(true, true)
	td = testDataIsValid{
		Name:     "mandatory fields with id",
		Actual:   nm,
		Expected: "true",
	}

	testCases = append(testCases, td)

	// all fields
	nm = m.Clone()
	nm.Attributes.ResetTestData()
	nm.Attributes.GetTestDataForStruct(false, true)
	td = testDataIsValid{
		Name:     "all fields",
		Actual:   nm,
		Expected: "true",
	}

	testCases = append(testCases, td)

	nm = m.Clone()
	nm.Attributes.ResetTestData()
	nm.Attributes.GetTestDataForStruct(false, false)
	td = testDataIsValid{
		Name:     "all fields without id",
		Actual:   nm,
		Expected: "true",
	}

	testCases = append(testCases, td)

	return testCases
}

func (m *model) GenerateTestData() *model {
	m.Attributes.ResetTestData()

	for _, v := range m.Attributes {
		_ = v.GetTestData()
	}

	return m
}

func (m *model) GetTestDataForStruct(mandatoryOnly bool, withID bool) *model {
	m.Attributes.ResetTestData()
	data := m.Clone()

	for _, v := range data.Attributes {
		if (mandatoryOnly && v.Nullable) || (!withID && v.Name == "ID") {
			continue
		}

		_ = v.GetTestData()
	}

	return data
}

func (m *model) GenerateTestDataForDuplicate(uniqueField *field) *model {
	b := m.Clone()
	b.GenerateTestData()
	for _, v := range b.Attributes {
		if v.Name == uniqueField.Name {
			f := m.Attributes.FindField(uniqueField.Name)
			if f != nil {
				v.TestData = f.TestData
			}
			break
		}
	}

	return b
}

func (m *model) GenerateTestDataForID() string {
	f := m.Attributes[0]
	d := f.GetTestDataForStruct(f.GetTestData())

	return strings.Replace(d, fieldNameID, strings.ToLower(fieldNameID), 1)
}

func (m *model) Clone() *model {
	c := &model{
		Name:     m.Name,
		rootPath: m.rootPath,
	}

	for _, attribute := range m.Attributes {
		a := &field{
			Name:         attribute.Name,
			PrimaryKey:   attribute.PrimaryKey,
			FieldType:    attribute.FieldType,
			DefaultValue: attribute.DefaultValue,
			MaxLength:    attribute.MaxLength,
			Nullable:     attribute.Nullable,
			Unique:       attribute.Unique,
			TestData:     attribute.TestData,
		}

		c.Attributes = append(c.Attributes, a)
	}

	return c
}

type testDataCRUD struct {
	Name        string
	Prepare     *model
	Actual      *model
	Expected    *model
	ExpectedErr string
}

type testDataIsValid struct {
	Name     string
	Actual   *model
	Expected string
}
