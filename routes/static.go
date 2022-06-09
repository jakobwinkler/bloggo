package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/jakobwinkler/bloggo/util"
)

func index(w http.ResponseWriter, r *http.Request) {
	const indexTemplate = "./templates/index.html.tmpl"
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		// This route will be hit if no others apply
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
		} else {
			err = util.ProcessTemplate(w, indexTemplate, nil)
			if err != nil {
				log.Printf("Error rendering template: %s", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		}
	}
}

func legal(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		io.WriteString(w, "Hello, legal!")
	}
}

func posts(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		const postsTemplate = "./templates/posts.html.tmpl"

		var data struct {
			Posts []postRoute
		}
		data.Posts = allPosts

		util.ProcessTemplate(w, postsTemplate, data)
		if err != nil {
			log.Printf("Error rendering template: %s", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

// TODO: is creating multiple routes (trailing slashes, ...) really a good idea?
func CreateStaticRoutes(mux *http.ServeMux, templateRoot string) {
	log.Println("Creating static routes...")

	mux.HandleFunc("/", index)
	mux.HandleFunc("/index.html", index)

	mux.HandleFunc("/legal", legal)
	mux.HandleFunc("/legal.html", legal)
	mux.HandleFunc("/legal/", legal)

	mux.HandleFunc("/posts", posts)
	mux.HandleFunc("/posts.html", posts)
	mux.HandleFunc("/posts/", posts)
}
