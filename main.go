package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jakobwinkler/bloggo/util"
)

func main() {
	log.Println("Launching Bloggo " + util.Version + " 🚀")

	// TODO: create static routes for index, legal, allposts (see current blog)
	// TODO: create dynamic routes for all files in posts/, but parse before and ignore if erroneous

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	// OK GO
	err := http.ListenAndServe(":8080", nil)

	// this should be unreachable
	if err != nil {
		log.Fatalf("HTTP server exited with error `%s`", err)
	}
}
