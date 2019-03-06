package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/oauth2"
)

type repo struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	Fork            bool   `json:"fork"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	ForksCount      int    `json:"forks_count"`
	Language        string `json:"language"`
	Archived        bool   `json:"archived"`

	Owner struct {
		Login     string `json:"login"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"owner"`
	License struct {
		Name string `json:"name"`
	} `json:"license"`
}

func main() {
	const base = "https://api.github.com"

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	url := base + "/user/repos?visibility=public&affiliation=owner,organization_member&sort=pushed"
	for url != "" {
		func() {
			res, err := tc.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			if res.StatusCode > 200 {
				log.Fatalf("call returned %d", res.StatusCode)
			}

			var repos []*repo
			if err := json.Unmarshal(b, &repos); err != nil {
				log.Fatal(err)
			}
			for _, r := range repos {
				if r.Private || r.Fork || r.Owner.Login == "splice" {
					continue
				}
				// TODO: if an organization repo, check that I'm a collaborator, otherwise skip
				// TODO: ignore some repos, e.g. my rsc.* clones, bug-related repos, etc.
				fmt.Println(r.FullName, r.Archived, r.Language, r.Description)
			}

			url = extractNextLink(res.Header)
		}()
	}
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
