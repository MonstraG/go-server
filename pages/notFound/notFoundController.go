package notFound

import (
	"go-server/helpers"
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

var notFoundTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notFound/notFound.gohtml"))
var notFoundPageData = pages.PageData{
	PageTitle: "404: page not found",
}

func GetHandler(w helpers.MyWriter, _ *http.Request) {
	err := notFoundTemplate.Execute(w, notFoundPageData)
	if err != nil {
		log.Fatal("Failed to render 404 page:\n", err)
	}
}
