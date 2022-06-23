package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jakobwinkler/bloggo/util"
)

var AllPosts []util.PostRoute

func post(w http.ResponseWriter, r *http.Request, path string) {
	const postTemplate = "./templates/post.tmpl.html"
	util.LogRequest(r)

	err := util.RefuseUnsupportedMethods(w, r)
	if err == nil {
		err, output, matter := util.RenderMarkdown(path)

		// Data required for template execution
		pageData := util.PageData{
			Posts:   AllPosts,
			Version: util.Version,
			Title:   matter.Title,
			Date:    matter.Date,
			Body:    template.HTML(output),
		}

		err = util.ProcessHTMLTemplate(w, postTemplate, pageData)
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
	for i, _ := range matches {
		m := matches[len(matches)-i-1]
		basename := filepath.Base(m)
		route := BaseRoute + strings.TrimSuffix(basename, filepath.Ext(basename))
		log.Printf("Routing %s via %s", m, route)

		err, matter := util.ParseFrontmatter(m)
		if err != nil {
			log.Fatalf("Error parsing frontmatter of %s: %s", m, err)
		}

		data := util.PostRoute{
			Route:  route,
			Matter: *matter,
		}
		AllPosts = append(AllPosts, data)
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			post(w, r, m)
		})
	}

	log.Printf("Created %d dynamic routes", len(matches))
}
