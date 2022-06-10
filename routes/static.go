package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jakobwinkler/bloggo/util"
)

func index(w http.ResponseWriter, r *http.Request) {
	// This route will be hit if no others apply
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		util.LogRequest(r)
		http.NotFound(w, r)
	} else {
		static(w, r, "static/index.md")
	}
}

func legal(w http.ResponseWriter, r *http.Request) {
	static(w, r, "static/legal.md")
}

func posts(w http.ResponseWriter, r *http.Request) {
	static(w, r, "static/posts.md")
}

func static(w http.ResponseWriter, r *http.Request, filepath string) {
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		const staticPageTemplate = "./templates/static.tmpl.html"

		err, output := util.RenderMarkdown(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageData := util.PageData{
			Posts:   AllPosts,
			Version: util.Version,
			Title:   filepath, // TODO: read front matter!
			Body:    template.HTML(output),
		}

		util.ProcessTemplate(w, staticPageTemplate, pageData)
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
