package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~mna/mna.dev/scripts/internal/types"
	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday/v2"
)

type website struct {
	Links       []*link
	IconCredits []*iconCredit

	// Posts, MicroPosts and Repos is set only when generating
	// the index.
	Posts      []*types.Post
	MicroPosts []*types.MicroPost
	Repos      []*types.Repo

	// CurrentPost is set only when generating specific pages,
	// in which case it is set to that page's MarkdownPost.
	CurrentPost *types.MarkdownPost
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

func main() {
	if len(os.Args) != 5 {
		log.Fatal("expect 4 arguments: posts, data, templates and destination directories")
	}

	posts, data, tpls, dst := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	dps, dms, drs, err := loadDataPostMicroRepo(data)
	if err != nil {
		log.Fatal(err)
	}
	lps, lms, lgs, err := loadLocalPostMicroPages(posts)
	if err != nil {
		log.Fatal(err)
	}

	t := parseTemplates(tpls)
	w := newWebsite()

	for _, post := range lps {
		if err := w.executePage(t, dst, post); err != nil {
			log.Fatal(err)
		}
		// once generated, merge with dps
		dps = append(dps, post.ToPost())
	}
	for _, page := range lgs {
		if err := w.executePage(t, dst, page); err != nil {
			log.Fatal(err)
		}
	}
	for _, micro := range lms {
		dms = append(dms, micro.ToMicroPost())
	}

	w.Posts = dps
	w.MicroPosts = dms
	w.Repos = drs

	// generate the index page
	if err := w.executeIndex(t, dst); err != nil {
		log.Fatal(err)
	}

	// TODO: add a SortedByDateDesc function on the
	// index to get a mixed list of all posts, micro-posts and repos by
	// published/updated date descending. This is what will be used in the
	// index to list them.
}

var funcs = template.FuncMap{
	"lower":    strings.ToLower,
	"markdown": toMarkdown,
	"markdownString": func(s string) template.HTML {
		return toMarkdown([]byte(s))
	},
}

func toMarkdown(b []byte) template.HTML {
	return template.HTML(blackfriday.Run(b))
}

func parseTemplates(dir string) *template.Template {
	t := template.New("root").Funcs(funcs)
	t, err := t.ParseGlob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func loadDataPostMicroRepo(dir string) (ps []*types.Post, ms []*types.MicroPost, rs []*types.Repo, err error) {
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
			return nil, nil, nil, err
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		for err == nil {
			v := newv()
			if err = dec.Decode(v); err != nil {
				if err == io.EOF {
					break
				}
				return nil, nil, nil, err
			}
			append(v)
		}
	}

	return ps, ms, rs, nil
}

func loadLocalPostMicroPages(dir string) (ps, ms, gs []*types.MarkdownPost, err error) {
	configs := make(map[string]*types.PostConfig)

	// first, walk to read the configuration TOML files
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		// stop if error walking dir
		if err != nil {
			return err
		}

		// this pass only cares about toml files
		ext := filepath.Ext(fi.Name())
		if ext != ".toml" {
			return nil
		}

		// extract the path relative to dir, without the extension
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = strings.TrimSuffix(rel, ext)

		var conf types.PostConfig
		_, err = toml.DecodeFile(path, &conf)
		if err != nil {
			return err
		}
		configs[rel] = &conf

		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	var posts []*types.MarkdownPost

	// next, walk to read the corresponding markdown files and collect
	// the list of all posts.
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		// stop if error walking dir
		if err != nil {
			return err
		}

		// this pass only cares about markdown files
		ext := filepath.Ext(fi.Name())
		if ext != ".md" {
			return nil
		}

		// extract the path relative to dir, without the extension
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = strings.TrimSuffix(rel, ext)

		// lookup the config for that post
		conf := configs[rel]
		if conf == nil {
			// ignore if there is no config
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		post := &types.MarkdownPost{
			Path:      rel,
			Title:     conf.Title,
			Published: conf.Published,
			Lead:      conf.Lead,
			Micro:     conf.Micro,
			Markdown:  b,
		}
		posts = append(posts, post)

		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// finally, split into posts, micro-posts and standalone pages
	for _, post := range posts {
		isPost := filepath.Dir(post.Path) == "posts"
		switch {
		case isPost && post.Micro:
			ms = append(ms, post)
		case isPost && !post.Micro:
			ps = append(ps, post)
		case !isPost && !post.Micro:
			// standalone pages cannot be micro-posts
			gs = append(gs, post)
		}
	}

	return ps, ms, gs, nil
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

func (w *website) executePage(t *template.Template, outDir string, post *types.MarkdownPost) error {
	f, err := os.Create(filepath.Join(outDir, post.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	w.CurrentPost = post
	if err := t.ExecuteTemplate(f, "post.html", w); err != nil {
		return err
	}
	return f.Close()
}

func (w *website) executeIndex(t *template.Template, outDir string) error {
	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	w.CurrentPost = nil
	if err := t.ExecuteTemplate(f, "index.html", w); err != nil {
		return err
	}
	return f.Close()
}
