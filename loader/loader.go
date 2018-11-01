package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	urlFormatCommits = "https://api.github.com/repos/hexdigest/gowrap/commits?path=templates/%s&per_page=1"
	urlFormatRaw     = "https://raw.githubusercontent.com/hexdigest/gowrap/%s/templates/%s"
	urlTree          = "https://api.github.com/repos/hexdigest/gowrap/git/trees/master?recursive=1"

	templatesPathPrefix = "templates/"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

//Loader fetches remote templates from the project's github page or any HTTPS URL
type Loader struct {
	client httpClient
}

//New returns Loader
func New(client httpClient) Loader {
	if client == nil {
		client = http.DefaultClient
	}

	return Loader{client: client}
}

//Load returns template contents and template URL or error
//path can be either any HTTPs URL or reference to the template on a project's github page
func (l Loader) Load(path string) (tmpl []byte, url string, err error) {
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		body, err := l.get(path)
		return body, path, err
	}
	if strings.HasPrefix(path, "file://") {
		body, err := ioutil.ReadFile(path[len("file://"):])
		return body, path, err
	}

	return l.fetchFromGithub(path)
}

//List returns a list of template names from the project's github page
func (l Loader) List() ([]string, error) {
	body, err := l.get(urlTree)
	if err != nil {
		return nil, err
	}

	tree := struct{ Tree []struct{ Path string } }{}

	if err := json.Unmarshal(body, &tree); err != nil {
		return nil, err
	}

	result := []string{}
	for _, leaf := range tree.Tree {
		if strings.HasPrefix(leaf.Path, templatesPathPrefix) && leaf.Path != templatesPathPrefix {
			path := strings.Replace(leaf.Path, templatesPathPrefix, "", 1)
			result = append(result, path)
		}
	}

	return result, nil
}

var errUnexpectedStatusCode = errors.New("unexpected status code")

func (l Loader) get(url string) (b []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errUnexpectedStatusCode, "%d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

var errTemplateNotFound = errors.New("remote template not found")

func (l Loader) fetchFromGithub(templateName string) ([]byte, string, error) {
	body, err := l.get(fmt.Sprintf(urlFormatCommits, templateName))
	if err != nil {
		return nil, "", err
	}

	commits := []struct{ SHA string }{}
	if err = json.Unmarshal(body, &commits); err != nil {
		return nil, "", errors.Wrap(err, "failed to decode commit info")
	}

	if len(commits) == 0 {
		return nil, "", errTemplateNotFound
	}

	url := fmt.Sprintf(urlFormatRaw, commits[0].SHA, templateName)

	contents, err := l.get(url)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to fetch template")
	}

	return contents, url, nil
}
