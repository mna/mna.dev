package sputnik

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
	"github.com/PuerkitoBio/goquery"
)

const (
	initialURL = "https://www.sputnikmusic.com/list.php?memberid=1142495"
	baseURL    = "https://www.sputnikmusic.com/"
)

type source struct {
}

func init() {
	datasource.Register("sputnik", &source{})
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

var rxTitleBestYear = regexp.MustCompile(`^Best\s.+\s(\d{4})$`)

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

	doc.Find("table.alt1 > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// skip, the "Lists" / "Create a list" header cells
			return
		}

		cells := s.Find("td")
		n := cells.Length()
		if n > 0 && n%2 != 0 {
			return
		}

		for i := 0; i < n/2; i++ {
			var published time.Time
			if dt := strings.TrimSpace(cells.Eq(i * 2).Text()); dt != "" {
				published, _ = time.Parse("01.02.06", dt)
			}
			anchor := cells.Eq(i*2 + 1).Find("a").First()
			link := anchor.AttrOr("href", "")
			if link != "" {
				link = baseURL + link
			}
			title := strings.TrimSpace(anchor.Text())
			if ms := rxTitleBestYear.FindStringSubmatch(title); ms != nil {
				// special-case for "Best of YYYY" lists, replace the published date
				// with the last day of the year of the best-of.
				year, _ := strconv.Atoi(ms[1])
				published = time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
			}

			post := &types.Post{
				URL:       link,
				Website:   "sputnikmusic",
				Title:     title,
				Published: published,
			}
			post.SetTags()
			emit <- post
		}
	})

	return "", nil
}
