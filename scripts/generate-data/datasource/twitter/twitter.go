package twitter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
	"github.com/PuerkitoBio/goquery"
)

const (
	initialURL = "https://twitter.com/___mna___/media"
	baseURL    = "https://twitter.com"
)

type source struct {
}

func init() {
	datasource.Register("twitter", &source{})
}

func (s *source) Generate(emit chan<- interface{}) error {
	cli := &http.Client{}

	var err error

	url := initialURL
	for url != "" {
		url, err = s.processPage(cli, url, emit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *source) processPage(client *http.Client, url string, emit chan<- interface{}) (next string, err error) {
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode > 200 {
		return "", fmt.Errorf("http status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	doc.Find(".content").Each(func(i int, s *goquery.Selection) {
		// ignore replies
		if s.Find(".ReplyingToContextBelowAuthor").Length() > 0 {
			return
		}

		var html string

		text := strings.TrimSpace(s.Find(".js-tweet-text-container").Text())
		link := s.Find(".time a").AttrOr("href", "")
		if link != "" {
			link = baseURL + link
			html, _ = generateEmbed(client, link)
		}

		var published time.Time
		dt := s.Find(".time ._timestamp").AttrOr("data-time", "")
		if dt != "" {
			epoch, err := strconv.ParseInt(dt, 10, 64)
			if err == nil {
				published = time.Unix(epoch, 0)
			}
		}

		post := &types.MicroPost{
			URL:       link,
			Website:   "twitter",
			Text:      text,
			RawHTML:   template.HTML(html),
			Published: published,
		}
		post.SetTags()
		emit <- post
	})
	return "", nil
}

func generateEmbed(cli *http.Client, url string) (string, error) {
	res, err := cli.Get(fmt.Sprintf("https://publish.twitter.com/oembed?url=%s", url))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode > 200 {
		return "", fmt.Errorf("http status code: %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var embed struct {
		HTML string `json:"html"`
	}
	if err := json.Unmarshal(b, &embed); err != nil {
		return "", err
	}
	return embed.HTML, nil
}
