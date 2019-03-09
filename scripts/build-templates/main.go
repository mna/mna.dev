package main

import (
	"fmt"
	"log"
	"os"
)

type website struct {
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("expect 3 arguments: posts directory, templates directory, destination directory")
	}

	posts, src, dst := os.Args[1], os.Args[2], os.Args[3]
	fmt.Println(">>>> ", posts, src, dst)
}
