package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

func main() {
	portFlag := flag.Int("port", 9000, "Port to listen on.")
	flag.Parse()

	log.Printf("listening on port %d...", *portFlag)

	// TODO: handle pretty URLs, so that /about serves /about.html, the same
	// way it will get done once deployed on netlify:
	// https://www.netlify.com/docs/redirects/#trailing-slash
	// Hmm this is a paid feature, maybe just generate the page without .html
	gz := gziphandler.MustNewGzipLevelHandler(gzip.DefaultCompression)
	h := gz(http.FileServer(http.Dir("public")))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), h); err != nil {
		log.Fatal(err)
	}
}
