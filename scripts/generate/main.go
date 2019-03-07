package main

import (
	"log"
	"os"

	"git.sr.ht/~mna/mna.dev/scripts/generate/datasource"
	_ "git.sr.ht/~mna/mna.dev/scripts/generate/datasource/github"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("want 1 argument, got %v", os.Args)
	}
	if err := datasource.Generate(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
