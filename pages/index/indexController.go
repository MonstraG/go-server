package index

import (
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/index/index.gohtml"))
var indexPageData = pages.PageData{
	PageTitle: "My todo list",
}

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Fatal("Failed to render index page", err)
	}
}
