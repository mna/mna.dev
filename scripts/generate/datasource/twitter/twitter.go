package twitter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate/datasource"
	"github.com/PuerkitoBio/goquery"
)

// TODO: can generate embed code by calling:
// https://publish.twitter.com/oembed?url=<the twitter url>
// and grabbing the html field on the returned JSON object.

const (
	initialURL = "https://twitter.com/___mna___/media"
	baseURL    = "https://twitter.com"
)

type post struct {
	URL       string
	Text      string
	Published time.Time
}

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

		text := strings.TrimSpace(s.Find(".js-tweet-text-container").Text())
		link := s.Find(".time a").AttrOr("href", "")
		if link != "" {
			link = baseURL + link
		}

		var published time.Time
		dt := s.Find(".time ._timestamp").AttrOr("data-time", "")
		if dt != "" {
			epoch, err := strconv.ParseInt(dt, 10, 64)
			if err == nil {
				published = time.Unix(epoch, 0)
			}
		}
		emit <- &post{
			URL:       link,
			Text:      text,
			Published: published,
		}
	})
	return "", nil
}
