package util

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func LogRequest(r *http.Request) {
	log.Printf("Got request: %s %s", r.Method, r.URL.String())
}

func ProcessTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	// Parse template
	temp, err := template.ParseFiles(templatePath)
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
