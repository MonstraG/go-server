package notFound

import (
	"html/template"
	"log"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("pages/notFound/notFound.gohtml"))

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		log.Fatal("Failed to render 404 page", err)
	}
}
