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

var todosTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/todos/todos.gohtml"))
var todosPageData = pages.PageData{
	PageTitle: "My todo list",
}

func (controller Controller) GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := todosTemplate.Execute(w, todosPageData)
	if err != nil {
		log.Fatal("Failed to render todos template", err)
	}
}

var todosListTemplate = template.Must(template.ParseFiles("pages/todos/todosList.gohtml"))

type ListDTO struct {
	Todos []Todo
}

func (controller Controller) GetListHandler(w http.ResponseWriter, _ *http.Request) {
	todos := controller.service.readTodos()

	pageModel := ListDTO{
		Todos: *todos,
	}

	err := todosListTemplate.Execute(w, pageModel)
	if err != nil {
		log.Fatal("Failed to render todos template", err)
	}
}

func (controller Controller) ApiGetHandler(w http.ResponseWriter, _ *http.Request) {
	todos := controller.service.readTodos()
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Print("Failed to marshal todos", err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Print("Failed to write todos", err)
	}
}

func (controller Controller) ApiPutHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	done := r.Form.Get("done") == "on"

	err = controller.service.setTodoDoneState(id, done)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
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

func (controller Controller) ApiDelHandler(w http.ResponseWriter, r *http.Request) {
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
