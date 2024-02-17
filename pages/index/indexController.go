package index

import (
	"html/template"
	"net/http"
)

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/index/index.gohtml"))
	var _ = indexTemplate.Execute(w, nil)
}
