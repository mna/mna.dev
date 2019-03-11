package main

import (
	"log"
	"os"

	"git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/bitbucket"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/github"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/gitlab"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/hypermegatop"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/sputnik"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/srht"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/staticpost"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/twitter"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate-data/datasource/zerovalue"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("want 1 argument, the destination directory, got %v", os.Args)
	}
	if err := datasource.Generate(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
