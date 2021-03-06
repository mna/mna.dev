package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
)

const (
	baseURL = "https://gitlab.com/api/v4/"
)

type repo struct {
	Description       string    `json:"description"`
	Visibility        string    `json:"visibility"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	HTTPURL           string    `json:"http_url_to_repo"`
	Archived          bool      `json:"archived"`
	ForksCount        int       `json:"forks_count"`
	StarCount         int       `json:"star_count"`
	CreatedAt         time.Time `json:"created_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	ForkedFromProject *struct {
		ID int `json:"id"`
	} `json:"forked_from_project"`
}

type source struct {
	base  string
	token string
}

func init() {
	tok := os.Getenv("GITLAB_API_TOKEN")
	if tok == "" {
		return
	}

	datasource.Register("gitlab", &source{
		base:  baseURL,
		token: os.Getenv("GITLAB_API_TOKEN"),
	})
}

func (s *source) Generate(emit chan<- interface{}) error {
	cli := &http.Client{}

	var err error

	url := s.base + "users/___mna___/projects?visibility=public&order_by=updated_at"
	for url != "" {
		url, err = s.processPage(cli, url, emit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *source) processPage(client *http.Client, url string, emit chan<- interface{}) (next string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// set the token
	req.Header.Set("Private-Token", s.token)

	// make the call
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 200 {
		return "", fmt.Errorf("http status code: %d", res.StatusCode)
	}

	var repos []*repo
	if err := json.Unmarshal(b, &repos); err != nil {
		return "", err
	}
	for _, r := range repos {
		if r.Visibility != "public" || r.ForkedFromProject != nil {
			continue
		}
		repo := &types.Repo{
			URL:         r.HTTPURL,
			Host:        "gitlab",
			Name:        r.NameWithNamespace,
			Description: r.Description,
			Created:     r.CreatedAt,
			Updated:     r.LastActivityAt,
			Stars:       r.StarCount,
			Forks:       r.ForksCount,
		}
		repo.SetTags()
		emit <- repo
	}

	url = extractNextLink(res.Header)
	return url, nil
}

var rxLinks = regexp.MustCompile(`<(.+?)>\s*;\s*rel="(\w+)"`)

func extractNextLink(h http.Header) string {
	links := h.Get("Link")
	ms := rxLinks.FindAllStringSubmatch(links, -1)
	for _, mms := range ms {
		if mms[2] == "next" {
			return mms[1]
		}
	}
	return ""
}
