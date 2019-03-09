package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
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

var funcs = template.FuncMap{
	"lower": strings.ToLower,
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("expect 3 arguments: posts directory, templates directory, destination directory")
	}

	posts, src, dst := os.Args[1], os.Args[2], os.Args[3]

	t := template.New("root").Funcs(funcs)
	t, err := t.ParseGlob(filepath.Join(src, "*"))
	if err != nil {
		log.Fatal(err)
	}

	w := &website{
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
			{Icon: "Stack Overflow", Name: "Freepik", AuthorURL: "https://www.freepik.com/", Website: "flaticon.com", WebsiteURL: "https://www.flaticon.com/", License: "CC 3.0 BY", LicenseURL: "http://creativecommons.org/licenses/by/3.0/"},
		},
	}
	if err := w.execute(t, "index.html", dst); err != nil {
		log.Fatal(err)
	}
	_ = posts
}

func (w *website) execute(t *template.Template, name string, outDir string) error {
	f, err := os.Create(filepath.Join(outDir, name))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.ExecuteTemplate(f, name, w); err != nil {
		return err
	}
	return f.Close()
}
