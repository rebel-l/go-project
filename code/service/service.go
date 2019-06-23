package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/config"
	"github.com/rebel-l/go-utils/osutils"
)

// Parameters defines parameters used for the go templates
type Parameters struct {
	Config   config.Config
	Packages []string
}

// NewParameters returns a new struct of Parameters prefilled by a config and the definition of packages
func NewParameters(cfg config.Config) Parameters {
	return Parameters{
		Config:   cfg,
		Packages: GetPackages().Get(),
	}
}

// Create the basic files for a service
func Create(projectPath string, params Parameters, commit git.CallbackAddAndCommit) error {
	pattern := filepath.Join("./code/service/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	var files []string
	for _, v := range getTemplateNames() {
		if err := ensurePath(projectPath, v.Name); err != nil {
			return err
		}

		filename := v.toFilenName()
		filename = filepath.Join(projectPath, filename)

		files = append(files, filename)
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create service main file: %s", err)
		}
		defer func() {
			_ = file.Close()
		}()

		if err = tmpl.ExecuteTemplate(file, v.Name, params); err != nil {
			return err
		}
	}

	return commit(files, "added go base file for service")
}

func ensurePath(projectPath, templateName string) error {
	parts := strings.Split(templateName, ".")
	path := projectPath
	if len(parts) > 1 {
		p := []string{projectPath}
		p = append(p, parts[:len(parts)-1]...)
		path = filepath.Join(p...)
	}

	if !osutils.FileOrPathExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func getTemplateNames() []templateName {
	return []templateName{
		{Name: "endpoint.doc.doc", FileExtension: "go"},
		{Name: "endpoint.doc.package", FileExtension: "go"},
		{Name: "endpoint.doc.swagger", FileExtension: "yml"},
		{Name: "endpoint.doc.web.index", FileExtension: "html"},
		{Name: "endpoint.ping.package", FileExtension: "go"},
		{Name: "endpoint.ping.ping", FileExtension: "go"},
		{Name: "main", FileExtension: "go"},
		{Name: "service.package", FileExtension: "go"},
		{Name: "service.service", FileExtension: "go"},
		{Name: "service.service_test", FileExtension: "go"},
	}
}

/*
TODO:
3. docs endpoint FIXME: index.html: replace title, headline, project description, author, url & license information
5. test file for ping endpoint
6. test file for docs endpoint
7. swagger definition FIXME: swagger.yml: license prefix & license fields in definition
8. later: auth client - permission request
9. investigate http.Server options
10. graceful service (see gorilla/mux)
11. middleware? ==> maybe service package
12. generate request uuid and add to logging: service:uuid
*/
