package srht

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
	"golang.org/x/oauth2"
)

const (
	baseURL    = "https://git.sr.ht"
	myUsername = "~mna"
)

type repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Visibility  string    `json:"visibility"`
}

type response struct {
	Errors  []map[string]string `json:"errors"`
	Next    *int                `json:"next"`
	Results []*repo             `json:"results"`
}

type source struct {
	base  string
	token string
}

func init() {
	tok := os.Getenv("SRHT_API_TOKEN")
	if tok == "" {
		return
	}

	datasource.Register("srht", &source{
		base:  baseURL,
		token: os.Getenv("SRHT_API_TOKEN"),
	})
}

func (s *source) Generate(emit chan<- interface{}) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.token, TokenType: "token"},
	)
	tc := oauth2.NewClient(ctx, ts)

	var err error

	url := s.base + "/api/repos"
	for url != "" {
		url, err = s.processPage(tc, url, emit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *source) processPage(client *http.Client, u string, emit chan<- interface{}) (next string, err error) {
	// make the call
	res, err := client.Get(u)
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

	var resp response
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	if len(resp.Errors) > 0 {
		return "", fmt.Errorf("%d error(s); first error: %#v", len(resp.Errors), resp.Errors[0])
	}

	for _, r := range resp.Results {
		if r.Visibility != "public" {
			continue
		}
		repo := &types.Repo{
			URL:         fmt.Sprintf("%s/%s/%s", baseURL, myUsername, r.Name),
			Host:        "sourcehut",
			Name:        r.Name,
			Description: r.Description,
			Created:     r.Created,
			Updated:     r.Updated,
		}
		repo.SetTags()
		emit <- repo
	}

	if resp.Next != nil && *resp.Next > 0 {
		parsed, err := url.Parse(u)
		if err != nil {
			return "", err
		}
		parsed.RawQuery = fmt.Sprintf("start=%d", *resp.Next)
		return parsed.String(), nil
	}
	return "", nil
}
