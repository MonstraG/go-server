package todos

import (
	"html/template"
	"net/http"
)

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/todos/todos.gohtml"))

	var todos = readTodos()

	var pageModel = DTO{
		Todos: *todos,
	}

	var _ = indexTemplate.Execute(w, pageModel)
}

func ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	id := parseIdPathValueSendErr(w, r)
	if id == 0 {
		return
	}
	err := parseFormSendErr(w, r)
	if err != nil {
		return
	}

	var done = r.Form.Get("done") == "on"

	var todos = readTodos()

	_, todo := findTodoById(todos, id)
	todo.Done = done

	writeTodos(todos)

	w.WriteHeader(200)
}

func ApiPutHandler(w http.ResponseWriter, r *http.Request) {
	err := parseFormSendErr(w, r)
	if err != nil {
		return
	}

	var title = r.Form.Get("title")
	if title == "" {
		w.WriteHeader(400)
		return
	}

	var todos = readTodos()

	*todos = append(*todos, Todo{
		Id:    generateNextId(todos),
		Title: title,
	})

	writeTodos(todos)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(200)
}

func ApiDelHandler(w http.ResponseWriter, r *http.Request) {
	id := parseIdPathValueSendErr(w, r)
	if id == 0 {
		return
	}

	updateTodos(deleteTodoAtIdAction(w, id))

	w.WriteHeader(200)
}
