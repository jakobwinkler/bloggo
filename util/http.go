package util

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
)

type PostRoute struct {
	Title string
	Route string
}

type PageData struct {
	Posts   []PostRoute
	Version string
	Title   string
	Body    template.HTML
}

func LogRequest(r *http.Request) {
	log.Printf("Got request: %s %s", r.Method, r.URL.String())
}

func ProcessTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	const masterTemplate = "./templates/master.tmpl.html"
	// Parse template
	temp, err := template.ParseFiles(masterTemplate, templatePath)
	if err != nil {
		return err
	}

	// Render template to temporary buffer so we can handle errors gracefully
	var buf bytes.Buffer
	err = temp.Execute(&buf, data)
	if err != nil {
		return err
	}

	// Finally, write the rendered template to the response
	buf.WriteTo(w)
	return nil
}

func RefuseUnsupportedMethods(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New(r.Method)
	}
	return nil
}

func RenderMarkdown(path string) (error, []byte) {
	// Read post contents
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading post content: %s", err)
		return err, nil
	}

	// Render HTML
	return nil, markdown.ToHTML(content, nil, nil)
}
