package hypermegatop

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"github.com/PuerkitoBio/goquery"
)

const baseURL = "http://hypermegatop.github.io"

type post struct {
	URL       string
	Title     string
	Lead      string
	Published time.Time
}

type source struct {
}

func init() {
	datasource.Register("hypermegatop", &source{})
}

func (s *source) Generate(emit chan<- interface{}) error {
	cli := &http.Client{}

	var err error

	url := baseURL
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

	doc.Find("#post-list li").Each(func(i int, s *goquery.Selection) {
		var published time.Time
		dt := s.Find("time.published").AttrOr("datetime", "")
		if dt != "" {
			published, _ = time.Parse("2006-01-02", dt)
		}
		anchor := s.Find("h2 a")
		link := anchor.AttrOr("href", "")
		if link != "" {
			link = baseURL + link
		}
		title := strings.TrimSpace(anchor.Text())
		lead := strings.TrimSpace(s.Find("p.abstract").Text())

		emit <- post{
			URL:       link,
			Title:     title,
			Lead:      lead,
			Published: published,
		}
	})

	next = doc.Find(".pager .previous a").AttrOr("href", "")
	if next != "" {
		next = baseURL + next
	}
	return next, nil
}
