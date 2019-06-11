package golang

// Import represents the information of a go import
type Import struct {
	Name  string
	Alias string
}

// Get returns the import command for go. If an alias is set, it returns it with alias, otherwise only the name of the package
func (i Import) Get() string {
	if i.Alias == "" {
		return "\"" + i.Name + "\""
	}

	return i.Alias + " \"" + i.Name + "\""
}

// Imports is a collection of Import
type Imports []Import

// GetNames returns a list of package names
func (i Imports) GetNames() []string {
	var names []string
	for _, v := range i {
		names = append(names, v.Name)
	}
	return names
}

// Get returns the list of import commands for go as it is provided by Import.Get()
func (i Imports) Get() []string {
	var imports []string
	for _, v := range i {
		imports = append(imports, v.Get())
	}
	return imports
}
