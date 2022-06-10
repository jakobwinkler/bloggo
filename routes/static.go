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
		static(w, r, "static/index.md", "templates/index.tmpl.html")
	}
}

func legal(w http.ResponseWriter, r *http.Request) {
	static(w, r, "static/legal.md", "templates/static.tmpl.html")
}

func posts(w http.ResponseWriter, r *http.Request) {
	static(w, r, "static/posts.md", "templates/posts.tmpl.html")
}

func static(w http.ResponseWriter, r *http.Request, filePath string, templatePath string) {
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		err, output, matter := util.RenderMarkdown(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageData := util.PageData{
			Posts:   AllPosts,
			Version: util.Version,
			Title:   matter.Title,
			Date:    matter.Date,
			Body:    template.HTML(output),
		}

		util.ProcessHTMLTemplate(w, templatePath, pageData)
		if err != nil {
			log.Printf("Error rendering template: %s", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

func rss(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		var data struct {
			Host        string
			Description string
			Posts       []util.PostRoute
		}

		// TODO: make this configurable
		data.Host = "blog.jwinkler.me"
		data.Description = "A blog"
		data.Posts = AllPosts

		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		util.ProcessTextTemplate(w, "templates/rss.tmpl.xml", data)
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

	mux.HandleFunc("/rss.xml", rss)
}
