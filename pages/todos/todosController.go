package todos

import (
	"encoding/json"
	"go-server/helpers"
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

	pageModel := ListDTO{
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
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	done := r.Form.Get("done") == "on"

	err = setTodoDoneState(id, done)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ApiPutHandler(w http.ResponseWriter, r *http.Request) {
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addTodo(title)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(http.StatusOK)
}

func ApiDelHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	err := deleteTodoById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
