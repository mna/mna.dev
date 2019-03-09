package zerovalue

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL    = "https://www.0value.com/"
	initialURL = "https://www.0value.com/build-a-blog-engine-in-go"
)

type post struct {
	URL       string
	Title     string
	Published time.Time
}

type source struct {
}

func init() {
	datasource.Register("0value", &source{})
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

	nav := doc.Find("nav")
	title := strings.TrimSpace(nav.Find(".middle").Text())
	next = nav.Find(".right a").AttrOr("href", "")
	if next != "" {
		next = baseURL + next
	}
	var published time.Time
	dt := doc.Find("main .meta time").AttrOr("datetime", "")
	if dt != "" {
		published, _ = time.Parse("2006-01-02", dt)
	}
	emit <- &post{
		URL:       url,
		Title:     title,
		Published: published,
	}

	return next, nil
}
