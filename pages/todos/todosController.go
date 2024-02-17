package todos

import (
	"encoding/json"
	"go-server/query"
	"html/template"
	"log"
	"net/http"
)

var todosTemplate = template.Must(template.ParseFiles("pages/todos/todos.gohtml"))

type ListDTO struct {
	Todos []Todo
}

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	todos := readTodos()

	var pageModel = ListDTO{
		Todos: *todos,
	}

	err := todosTemplate.Execute(w, pageModel)
	if err != nil {
		log.Fatal("Failed to render todos template", err)
	}
}

func ApiGetHandler(w http.ResponseWriter, _ *http.Request) {
	todos := readTodos()
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Print("Failed to marshal todos", err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Print("Failed to write todos", err)
	}
}

func ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	id := query.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	err := query.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	var done = r.Form.Get("done") == "on"

	runTodosAction(changeDoneAction(w, id, done))

	w.WriteHeader(200)
}

func ApiPutHandler(w http.ResponseWriter, r *http.Request) {
	err := query.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	var title = r.Form.Get("title")
	if title == "" {
		w.WriteHeader(400)
		return
	}

	runTodosAction(addTodoAction(title))

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(200)
}

func ApiDelHandler(w http.ResponseWriter, r *http.Request) {
	id := query.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	runTodosAction(deleteTodoAtIdAction(w, id))

	w.WriteHeader(200)
}
