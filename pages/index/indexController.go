package index

import (
	"html/template"
	"log"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("pages/index/index.gohtml"))

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal("Failed to render index page", err)
	}
}
