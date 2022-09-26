package pages

import (
	"net/http"
	"stl/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func SitemapXSL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xsl")
	renderHTML(w, "sitemap.xsl", nil)
	return
}

func SitemapIndex(w http.ResponseWriter, r *http.Request) {
	xml, err := model.SitemapIndex()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(xml))
	return
}

func SitemapPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	part := vars["part"]
	idPart := strings.TrimSuffix(strings.TrimPrefix(part, "part-"), ".xml")
	numPart, err := strconv.Atoi(idPart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if numPart <= 0 {
		http.Error(w, "Part number must be > 0", http.StatusBadRequest)
		return
	}

	xml, err := model.SitemapPart(numPart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(xml))
	return
}
