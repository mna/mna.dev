package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"git.sr.ht/~mna/mna.dev/scripts/generate/datasource"
	"golang.org/x/oauth2"
)

const (
	githubBaseURL = "https://api.github.com"
	myUsername    = "mna"
)

type repo struct {
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	Private         bool     `json:"private"`
	HTMLURL         string   `json:"html_url"`
	Description     string   `json:"description"`
	Fork            bool     `json:"fork"`
	StargazersCount int      `json:"stargazers_count"`
	WatchersCount   int      `json:"watchers_count"`
	ForksCount      int      `json:"forks_count"`
	Language        string   `json:"language"`
	Archived        bool     `json:"archived"`
	Topics          []string `json:"topics"`

	Owner struct {
		Login     string `json:"login"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"owner"`

	License struct {
		Name string `json:"name"`
	} `json:"license"`
}

func (r *repo) ignoredTopic() bool {
	for _, topic := range r.Topics {
		if topic == "ignore" {
			return true
		}
	}
	return false
}

type source struct {
	base  string
	token string
}

func init() {
	tok := os.Getenv("GITHUB_API_TOKEN")
	if tok == "" {
		return
	}

	datasource.Register("github", &source{
		base:  githubBaseURL,
		token: os.Getenv("GITHUB_API_TOKEN"),
	})
}

func (s *source) Generate(emit chan<- interface{}) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.token},
	)
	tc := oauth2.NewClient(ctx, ts)

	var err error

	url := s.base + "/user/repos?visibility=public&affiliation=owner,organization_member&sort=pushed"
	for url != "" {
		url, err = s.processPage(tc, url, emit)
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

	// set the media type to receive topics
	req.Header.Set("Accept", "application/vnd.github.mercy-preview+json")

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
		if r.Private || r.Fork || r.Owner.Login == "splice" || r.ignoredTopic() {
			continue
		}
		if r.Owner.Type == "Organization" {
			// check that I'm a contributer, otherwise skip
			res, err := client.Get(fmt.Sprintf(s.base+"/repos/%s/%s/commits?author=%s", r.Owner.Login, r.Name, myUsername))
			if err != nil {
				return "", err
			}
			var array []interface{}
			json.NewDecoder(res.Body).Decode(&array)
			res.Body.Close()
			if res.StatusCode != 200 || len(array) == 0 {
				// not a contributor
				continue
			}
		}
		emit <- r
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
