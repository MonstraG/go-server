package index

import (
	"go-server/helpers"
	"go-server/pages"
	"html/template"
)

var indexTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/index/index.gohtml"))
var indexPageData = pages.PageData{
	PageTitle: "Homepage",
}

func GetHandler(w *helpers.MyWriter, _ *helpers.MyRequest) {
	w.ExecuteTemplate(indexTemplate, indexPageData)
}
