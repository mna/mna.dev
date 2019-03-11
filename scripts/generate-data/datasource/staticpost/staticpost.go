package staticpost

import (
	"time"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
)

type source struct {
}

func init() {
	datasource.Register("staticpost", &source{})
}

func (s *source) Generate(emit chan<- interface{}) error {
	emit <- &types.Post{
		Website:   "GopherAcademy",
		URL:       "https://blog.gopheracademy.com/advent-2014/goquery/",
		Title:     "goquery: a little like that j-thing",
		Published: time.Date(2014, 12, 12, 0, 0, 0, 0, time.UTC),
	}

	emit <- &types.Post{
		Website:   "Splice",
		URL:       "https://splice.com/blog/lesser-known-features-go-test/",
		Title:     "Lesser-known features of go test",
		Published: time.Date(2014, 9, 3, 0, 0, 0, 0, time.UTC),
	}

	emit <- &types.Post{
		Website:   "Splice",
		URL:       "https://splice.com/blog/going-extra-mile-golint-go-vet/",
		Title:     "Going the extra mile: golint and go vet",
		Published: time.Date(2014, 7, 10, 0, 0, 0, 0, time.UTC),
	}
	return nil
}
