package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
)

const (
	baseURL    = "https://%s:%s@api.bitbucket.org/2.0"
	myUsername = "___mna___"
)

type repo struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Language    string    `json:"language"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	IsPrivate   bool      `json:"is_private"`
	Description string    `json:"description"`
	Parent      *struct {
		UUID string `json:"uuid"`
	} `json:"parent"`
}

type response struct {
	Page   int     `json:"page"`
	Next   string  `json:"next"`
	Values []*repo `json:"values"`
}

type source struct {
	base  string
	token string
}

func init() {
	tok := os.Getenv("BITBUCKET_API_TOKEN")
	if tok == "" {
		return
	}

	datasource.Register("bitbucket", &source{
		base:  baseURL,
		token: os.Getenv("BITBUCKET_API_TOKEN"),
	})
}

func (s *source) Generate(emit chan<- interface{}) error {
	cli := &http.Client{}

	var err error

	url := fmt.Sprintf(s.base, myUsername, s.token) + fmt.Sprintf("/repositories/%s?q=is_private=false&sort=-updated_on", myUsername)
	for url != "" {
		url, err = s.processPage(cli, url, emit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *source) processPage(client *http.Client, u string, emit chan<- interface{}) (next string, err error) {
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
		return "", fmt.Errorf("http status code: %d (%s)", res.StatusCode, string(b))
	}

	var resp response
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	for _, r := range resp.Values {
		if r.IsPrivate || r.Parent != nil {
			continue
		}
		emit <- r
	}

	if resp.Next != "" {
		parsed, err := url.Parse(resp.Next)
		if err != nil {
			return "", err
		}
		parsed.User = url.UserPassword(myUsername, s.token)
		return parsed.String(), nil
	}
	return "", nil
}
