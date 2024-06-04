package todos

import (
	"encoding/json"
	"go-server/helpers"
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

var todosTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/todos/tmpl/todos.gohtml"))
var todosPageData = pages.PageData{
	PageTitle: "My todo list",
}

func (controller Controller) GetHandler(w helpers.MyWriter, _ *http.Request) {
	err := todosTemplate.Execute(w, todosPageData)
	if err != nil {
		log.Fatal("Failed to render todos template:\n", err)
	}
}

var todosListTemplate = template.Must(template.ParseFiles("pages/todos/tmpl/todosList.gohtml"))

type ListDTO struct {
	Todos []Todo
}

func (controller Controller) GetListHandler(w helpers.MyWriter, _ *http.Request) {
	todos := controller.service.readTodos()

	pageModel := ListDTO{
		Todos: *todos,
	}

	err := todosListTemplate.Execute(w, pageModel)
	if err != nil {
		log.Fatal("Failed to render todos template:\n", err)
	}
}

func (controller Controller) ApiGetHandler(w helpers.MyWriter, _ *http.Request) {
	todos := controller.service.readTodos()
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Print("Failed to marshal todos:\n", err)
	}
	w.WriteSilent(bytes)
}

func (controller Controller) ApiPutHandler(w helpers.MyWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	ok := helpers.ParseFormRespondErr(w, r)
	if !ok {
		return
	}

	done := r.Form.Get("done") == "on"

	err := controller.service.setTodoDoneState(id, done)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiPostHandler(w helpers.MyWriter, r *http.Request) {
	ok := helpers.ParseFormRespondErr(w, r)
	if !ok {
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	controller.service.addTodo(title)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiDelHandler(w helpers.MyWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	err := controller.service.deleteTodoById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
