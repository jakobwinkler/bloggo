package routes

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jakobwinkler/bloggo/util"
)

type postRoute struct {
	Title string
	Route string
}

var allPosts []postRoute

func post(w http.ResponseWriter, r *http.Request, path string) {
	const postTemplate = "./templates/post.html.tmpl"
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		// Read post contents
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Error reading post content: %s", err)
			http.Error(w, "Error reading post", http.StatusInternalServerError)
			return
		}

		// Data required for template execution
		type postData struct {
			Title string
			Body  string
			Route string
		}

		data := postData{
			Title: path,
			Body:  string(content),
		}

		err = util.ProcessTemplate(w, postTemplate, data)
		if err != nil {
			log.Printf("Error rendering template: %s", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

func CreateDynamicRoutes(mux *http.ServeMux, templateRoot string, blogRoot string) {
	log.Println("Creating dynamic routes...")

	// Find all blog posts
	matches, err := filepath.Glob(blogRoot + "/*.md")
	if err != nil {
		log.Fatalf("Error accessing files in %s: %s", blogRoot, err)
	}

	// Create routes for all posts and store them so we can show a list
	const BaseRoute = "/posts/"
	for _, m := range matches {
		basename := filepath.Base(m)
		route := BaseRoute + strings.TrimSuffix(basename, filepath.Ext(basename))
		log.Printf("Routing %s via %s", m, route)

		data := postRoute{
			Title: "<TODO: parse front matter>",
			Route: route,
		}
		allPosts = append(allPosts, data)
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			post(w, r, m)
		})
	}

	log.Printf("Created %d dynamic routes", len(matches))
}
