package util

import (
	"bytes"
	"errors"
	htemplate "html/template"
	"io/ioutil"
	"log"
	"net/http"
	ttemplate "text/template"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type PostRoute struct {
	Route  string
	Matter FrontMatter
}

type PageData struct {
	Posts   []PostRoute
	Version string
	Title   string
	Date    string
	Body    htemplate.HTML
}

type FrontMatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

func LogRequest(r *http.Request) {
	log.Printf("Got request: %s %s", r.Method, r.URL.String())
}

func ProcessTextTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	// Parse template
	temp, err := ttemplate.ParseFiles(templatePath)
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

func ProcessHTMLTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	const masterTemplate = "./templates/master.tmpl.html"

	// Parse template
	temp, err := htemplate.ParseFiles(masterTemplate, templatePath)
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

func RenderMarkdown(path string) (error, []byte, *FrontMatter) {
	// Read post contents
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading post file: %s", err)
		return err, nil, nil
	}

	matter := FrontMatter{}
	rest, err := frontmatter.Parse(bytes.NewReader(content), &matter)
	if err != nil {
		log.Printf("Error reading front matter: %s", err)
		return err, nil, nil
	}

	// Render HTML
	extensions := parser.NoIntraEmphasis | parser.Tables |
		parser.Strikethrough | parser.SpaceHeadings | parser.Footnotes |
		parser.HeadingIDs | parser.AutoHeadingIDs | parser.DefinitionLists |
		parser.Attributes | parser.SuperSubscript | parser.Includes |
		parser.Mmark
	parser := parser.NewWithExtensions(extensions)
	return nil, markdown.ToHTML(rest, parser, nil), &matter
}

func ParseFrontmatter(path string) (error, *FrontMatter) {
	// Read post contents
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading post file: %s", err)
		return err, nil
	}

	matter := FrontMatter{}
	_, err = frontmatter.Parse(bytes.NewReader(content), &matter)
	if err != nil {
		log.Printf("Error reading front matter: %s", err)
		return err, nil
	}

	return nil, &matter
}
