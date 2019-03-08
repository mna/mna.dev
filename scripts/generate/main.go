package main

import (
	"log"
	"os"

	"git.sr.ht/~mna/mna.dev/scripts/generate/datasource"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/bitbucket"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/github"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/gitlab"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/hypermegatop"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/srht"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/twitter"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/zerovalue"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("want 1 argument, got %v", os.Args)
	}
	if err := datasource.Generate(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
