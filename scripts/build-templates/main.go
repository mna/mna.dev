package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type website struct {
	Links []*link
}

type link struct {
	Website  string
	Username string
	URL      string
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
			{"Twitter", "___mna___", "https://twitter.com/___mna___/"},
			{"GitHub", "mna", "https://github.com/mna"},
			{"StackOverflow", "mna", "https://stackoverflow.com/users/1094941/mna"},
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
