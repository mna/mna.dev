package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
)

type website struct {
	Links       []*link
	IconCredits []*iconCredit
}

type link struct {
	Website  string
	Username string
	URL      string
}

type iconCredit struct {
	Icon       string
	Name       string
	AuthorURL  string
	Website    string
	WebsiteURL string
	License    string
	LicenseURL string
}

type index struct {
	Website    *website
	Posts      []*types.Post
	MicroPosts []*types.MicroPost
	Repos      []*types.Repo
}

func main() {
	if len(os.Args) != 5 {
		log.Fatal("expect 4 arguments: posts, data, templates and destination directories")
	}

	posts, data, tpls, dst := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	t := parseTemplates(tpls)
	ps, ms, rs := loadPostMicroRepo(data)
	i := &index{
		Website:    newWebsite(),
		Posts:      ps,
		MicroPosts: ms,
		Repos:      rs,
	}
	if err := i.execute(t, dst); err != nil {
		log.Fatal(err)
	}
	_ = posts
}

var funcs = template.FuncMap{
	"lower": strings.ToLower,
}

func parseTemplates(dir string) *template.Template {
	t := template.New("root").Funcs(funcs)
	t, err := t.ParseGlob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func loadPostMicroRepo(dir string) (ps []*types.Post, ms []*types.MicroPost, rs []*types.Repo) {
	append := func(v interface{}) {
		switch v := v.(type) {
		case *types.Post:
			ps = append(ps, v)
		case *types.MicroPost:
			ms = append(ms, v)
		case *types.Repo:
			rs = append(rs, v)
		}
	}

	files := map[string]func() interface{}{
		"post.json":  func() interface{} { return new(types.Post) },
		"mpost.json": func() interface{} { return new(types.MicroPost) },
		"repo.json":  func() interface{} { return new(types.Repo) },
	}

	for fnm, newv := range files {
		f, err := os.Open(filepath.Join(dir, fnm))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		for err == nil {
			v := newv()
			if err = dec.Decode(v); err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			append(v)
		}
	}

	return ps, ms, rs
}

func newWebsite() *website {
	return &website{
		Links: []*link{
			{"About", "", "/about"},
			{"Twitter", "___mna___", "https://twitter.com/___mna___/"},
			{"GitHub", "mna", "https://github.com/mna"},
			{"StackOverflow", "mna", "https://stackoverflow.com/users/1094941/mna"},
		},
		IconCredits: []*iconCredit{
			{Icon: "About", Name: "Dave Gandy", AuthorURL: "https://www.flaticon.com/authors/dave-gandy", Website: "flaticon.com", WebsiteURL: "https://www.flaticon.com/", License: "CC 3.0 BY", LicenseURL: "http://creativecommons.org/licenses/by/3.0/"},
			{Icon: "GitHub", Name: "Dave Gandy", AuthorURL: "https://www.flaticon.com/authors/dave-gandy", Website: "flaticon.com", WebsiteURL: "https://www.flaticon.com/", License: "CC 3.0 BY", LicenseURL: "http://creativecommons.org/licenses/by/3.0/"},
			{Icon: "Twitter", Name: "Katarina Stefanikova", AuthorURL: "https://www.flaticon.com/authors/katarina-stefanikova", Website: "flaticon.com", WebsiteURL: "https://www.flaticon.com/", License: "CC 3.0 BY", LicenseURL: "http://creativecommons.org/licenses/by/3.0/"},
			{Icon: "Stack Overflow", Name: "Pixel perfect", AuthorURL: "https://www.flaticon.com/authors/pixel-perfect", Website: "flaticon.com", WebsiteURL: "https://www.flaticon.com/", License: "CC 3.0 BY", LicenseURL: "http://creativecommons.org/licenses/by/3.0/"},
		},
	}
}

func (i *index) execute(t *template.Template, outDir string) error {
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.ExecuteTemplate(f, "index.html", i); err != nil {
		return err
	}
	return f.Close()
}
