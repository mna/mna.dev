package srht

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"golang.org/x/oauth2"
)

const (
	baseURL = "https://git.sr.ht"
)

type repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Visibility  string    `json:"visibility"`
}

type response struct {
	Next    int     `json:"next"`
	Results []*repo `json:"results"`
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

	for _, r := range resp.Results {
		if r.Visibility != "public" {
			continue
		}
		emit <- r
	}

	if resp.Next > 0 {
		// TODO: for now, the next page query doesn't seem to work on sr.ht (always
		// returns the same results regardless). Either the doc is wrong at
		// https://man.sr.ht/api-conventions.md#api-conventions or there is a bug.
		/*
			parsed, err := url.Parse(u)
			if err != nil {
				return "", err
			}
			parsed.RawQuery = fmt.Sprintf("get=%d", resp.Next)
			return parsed.String(), nil
		*/
	}
	return "", nil
}