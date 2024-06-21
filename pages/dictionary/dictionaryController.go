package dictionary

import (
	"go-server/helpers"
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/dictionary/dictionary.gohtml"))
var indexPageData = pages.PageData{
	PageTitle: "Dictionary",
}

func GetHandler(w helpers.MyWriter, _ *http.Request) {
	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Fatal("Failed to render index page:\n", err)
	}
}
