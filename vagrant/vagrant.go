package vagrant

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/rand"
	goutil "github.com/rebel-l/go-utils/strings"
)

const (
	fileNameBootstrap    = "vm/bootstrap.sh"
	fileNameVagrant      = "Vagrantfile"
	hostnamePattern      = "%s.test"
	ipPattern            = "192.168.%d.%d"
	min                  = 3
	max                  = 253
	templateKeyBootstrap = "bootstrap"
	templateKeyVagrant   = "vagrantfile"
)

var (
	params *Vagrant
	files  = map[string]string{
		templateKeyVagrant:   fileNameVagrant,
		templateKeyBootstrap: fileNameBootstrap,
	}
)

type Vagrant struct {
	ServiceName     string
	IP              string
	Hostname        string
	HostnameAliases []string
}

func newVagrant(project string, hostname string, domainPrefixes []string) *Vagrant {
	projectParts := strings.Split(project, "-")
	for k, v := range projectParts {
		projectParts[k] = strings.Title(v)
	}

	var hostnameAliases []string
	for _, v := range domainPrefixes {
		hostnameAliases = append(hostnameAliases, fmt.Sprintf("%s.%s", v, hostname))
	}

	return &Vagrant{
		ServiceName:     strings.Join(projectParts, ""),
		IP:              fmt.Sprintf(ipPattern, rand.Int(min, max), rand.Int(min, max)),
		Hostname:        hostname,
		HostnameAliases: hostnameAliases,
	}
}

func Prepare(project string) {
	if kind.Get() != kind.Service || !confirmation() {
		return
	}

	hostname := fmt.Sprintf(hostnamePattern, project)
	domainPrefixes := askForDomainPrefixes(hostname)
	params = newVagrant(project, hostname, domainPrefixes)
}

func Setup(path string, commit git.CallbackAddAndCommit, step int) error {
	if params == nil {
		return nil
	}

	pattern := filepath.Join("./vagrant/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	filenames, err := createFiles(path, tmpl)
	if err != nil {
		return err
	}

	return commit(filenames, "setup vagrant", step)
}

func createFiles(path string, tmpl *template.Template) ([]string, error) {
	var fileList []string

	for k, v := range files {
		filename, err := createFile(path, tmpl, k, v)
		if err != nil {
			return nil, err
		}

		fileList = append(fileList, filename)
	}

	return fileList, nil
}

func createFile(path string, tmpl *template.Template, tmplKey string, filename string) (string, error) {
	filename = filepath.Join(path, filename)
	if osutils.FileOrPathExists(filename) {
		return "", nil
	}

	subPath := filepath.Dir(filename)
	if path != subPath {
		if err := osutils.CreateDirectoryIfNotExists(subPath); err != nil {
			return "", err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, tmplKey, params); err != nil {
		return "", fmt.Errorf("failed to write template to file %s: %w", filename, err)
	}

	return filename, nil
}

func confirmation() bool {
	answer := "y"
	t := prompt.Input("Add Vagrant to this project? [Y/n] ", func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix([]prompt.Suggest{}, d.GetWordBeforeCursor(), true)
	})

	if t != "" {
		answer = t
	}

	return strings.ToLower(answer) == "y"
}

func askForDomainPrefixes(hostname string) []string {
	s := prompt.Input(
		fmt.Sprintf("Add a list of subdomains (comma seperated list as prefixes for %s): ", hostname),
		func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		},
	)

	return goutil.SplitTrimSpace(s, ",")
}
