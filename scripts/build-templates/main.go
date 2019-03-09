package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type website struct {
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("expect 3 arguments: posts directory, templates directory, destination directory")
	}

	posts, src, dst := os.Args[1], os.Args[2], os.Args[3]
	t, err := template.ParseGlob(filepath.Join(src, "*"))
	if err != nil {
		log.Fatal(err)
	}

	if err := execute(t, "index.html", dst); err != nil {
		log.Fatal(err)
	}
	_ = posts
}

func execute(t *template.Template, name string, outDir string) error {
	f, err := os.Create(filepath.Join(outDir, name))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.ExecuteTemplate(f, name, nil); err != nil {
		return err
	}
	return f.Close()
}
