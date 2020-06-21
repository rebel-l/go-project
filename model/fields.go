package model

import (
	"strings"

	"github.com/rebel-l/go-project/golang"
)

type fields []*field

func (f fields) GetSQLFieldNames() string {
	var fieldNames []string

	for _, v := range f {
		fieldNames = append(fieldNames, v.GetSQlFieldName())
	}

	return strings.Join(fieldNames, ", ")
}

func (f fields) GetSQLFieldNamesWithoutID() string {
	var fieldNames []string

	if len(f) < 2 {
		return ""
	}

	for _, v := range f[1:] {
		fieldNames = append(fieldNames, v.GetSQlFieldName())
	}

	return strings.Join(fieldNames, ", ")
}

func (f fields) GetPackages(projectRootPath string) ([]string, error) {
	var packages []string

	for _, v := range f {
		switch v.FieldType {
		case fieldTypeUUID:
			if err := golang.Get(projectRootPath, packageUUID); err != nil {
				return nil, err
			}

			packages = append(packages, packageUUID)
		}
	}

	return packages, nil
}

func (f fields) GetNotNullableFieldsWithComparison(receiver string) []string {
	var fields []string

	for _, v := range f {
		if !v.Nullable {
			fields = append(fields, v.GetEmptyComparison(receiver))
		}
	}

	return fields
}

func (f fields) GetTestData(mandatoryOnly bool, withID bool) string {
	var data []string
	for _, v := range f {
		if (mandatoryOnly && v.Nullable) || (!withID && v.Name == "ID") {
			continue
		}

		data = append(data, v.GetTestData())
	}

	return strings.Join(data, ", ")
}

func (f fields) CountMandatory() int {
	var i int
	for _, v := range f {
		if !v.Nullable {
			i++
		}
	}

	return i
}
