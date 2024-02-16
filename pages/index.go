package pages

import (
	"example-go-server/models"
	"html/template"
	"net/http"
)

type TodoPageData struct {
	PageTitle string
	Todos     []models.Todo
}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/index.html"))
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []models.Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	var _ = indexTemplate.Execute(w, data)
}
