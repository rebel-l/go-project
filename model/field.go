package model

type fields []field

type field struct {
	name         string
	fieldType    string
	defaultValue string
	nullable     bool
	maxLength    int
}
