package vagrant

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/c-bata/go-prompt"

	"github.com/rebel-l/go-project/kind"

	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/lib/print"
	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/go-utils/rand"
	goutil "github.com/rebel-l/go-utils/strings"
)

const (
	fileName        = "Vagrantfile"
	hostnamePattern = "%s.test"
	ipPattern       = "192.168.%d.%d"
	min             = 3
	max             = 253
	templateKey     = "vagrantfile"
)

var (
	params *Vagrant
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
	if !confirmation() {
		return
	}

	hostname := fmt.Sprintf(hostnamePattern, project)
	domainPrefixes := askForDomainPrefixes(hostname)
	params = newVagrant(project, hostname, domainPrefixes)
}

func Setup(path string, commit git.CallbackAddAndCommit, step int) error {
	if kind.Get() != kind.Service {
		return nil
	}

	if params == nil {
		return nil
	}

	filename := filepath.Join(path, fileName)
	if osutils.FileOrPathExists(filename) {
		print.Info("Skip creating a vagrant file as it already exists")
		return nil
	}

	pattern := filepath.Join("./vagrant/tmpl", "*.tmpl")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create vagrant file: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err = tmpl.ExecuteTemplate(file, templateKey, params); err != nil {
		return err
	}

	return commit([]string{filename}, "added Vagrantfile", step)
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
