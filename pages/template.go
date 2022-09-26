package pages

import (
	"html/template"
	"net/http"
)

var Templates *template.Template

func renderHTML(w http.ResponseWriter, templateName string, viewArgs map[string]interface{}) {
	if viewArgs == nil {
		viewArgs = map[string]interface{}{}
	}
	s1 := Templates.Lookup(templateName)
	if s1 == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	err := s1.ExecuteTemplate(w, templateName, viewArgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func renderText(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(text))
}

func renderJsError(w http.ResponseWriter, errorText string, status int) {
	s1 := Templates.Lookup("error-response.html")
	if s1 == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	err := s1.ExecuteTemplate(w, "error-response.html", errorText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
