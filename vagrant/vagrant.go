package vagrant

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/rebel-l/go-project/dialog"
	"github.com/rebel-l/go-project/git"
	"github.com/rebel-l/go-project/kind"
	"github.com/rebel-l/go-project/template"
	"github.com/rebel-l/go-utils/randutils"
	"github.com/rebel-l/go-utils/stringsutils"

	"github.com/c-bata/go-prompt"
)

const (
	hostnamePattern = "%s.test"
	ipPattern       = "192.168.%d.%d"
	min             = 3
	max             = 253

	templateKeyVagrant = "vagrantfile"

	templateKeyBootstrap = "bootstrap"

	templateKeyBashRC  = "bashrc"
	templateKeyProfile = "profile"

	templateNginxDHParam                    = "nginxDHParam"
	templateKeyNginxSiteConf                = "nginxSiteConf"
	templateKeyNginxSnippetsCertificateConf = "nginxSnippetsCertificateConf"
	templateKeyNginxSnippetsSSLParamsConf   = "nginxSnippetsSSLParamsConf"
	templateKeyNginxSnippetsTrailingSlash   = "nginxSnippetsTrailingSlashConf"
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

func (v *Vagrant) getFileConfig() map[string]string {
	pathVM := "vm"

	pathHome := path.Join(pathVM, "home", "vagrant")

	pathNginx := path.Join(pathVM, "etc", "nginx")
	pathNginxSiteConf := path.Join(pathNginx, "sites-available")
	pathNginxSnippets := path.Join(pathNginx, "snippets")

	return map[string]string{
		templateKeyVagrant:                      "Vagrantfile",
		templateKeyBootstrap:                    path.Join(pathVM, "bootstrap.sh"),
		templateKeyBashRC:                       path.Join(pathHome, ".bashrc"),
		templateKeyProfile:                      path.Join(pathHome, ".profile"),
		templateNginxDHParam:                    path.Join(pathNginx, "dhparam.pem"),
		templateKeyNginxSiteConf:                path.Join(pathNginxSiteConf, v.Hostname),
		templateKeyNginxSnippetsCertificateConf: path.Join(pathNginxSnippets, fmt.Sprintf("%s-certificate.conf", v.Hostname)),
		templateKeyNginxSnippetsSSLParamsConf:   path.Join(pathNginxSnippets, "ssl-params.conf"),
		templateKeyNginxSnippetsTrailingSlash:   path.Join(pathNginxSnippets, "trailingslash.conf"),
	}
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
		IP:              fmt.Sprintf(ipPattern, randutils.Int(min, max), randutils.Int(min, max)),
		Hostname:        hostname,
		HostnameAliases: hostnameAliases,
	}
}

func Prepare(project string) {
	if kind.Get() != kind.Service || !dialog.Confirmation("Add Vagrant to this project?") {
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

	filenames, err := template.CreateFilesWithTemplatePath(pattern, path, params.getFileConfig(), params)
	if err != nil {
		return fmt.Errorf("vagrant setup failed: %w", err)
	}

	return commit(filenames, "setup vagrant", step)
}

func askForDomainPrefixes(hostname string) []string {
	s := prompt.Input(
		fmt.Sprintf("Add a list of subdomains (comma seperated list as prefixes for %s): ", hostname),
		func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		},
	)

	return stringsutils.SplitTrimSpace(s, ",")
}
