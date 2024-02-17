package pages

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type IndexPageModel struct {
	PageTitle string
}

type TodosDTO struct {
	Todos []Todo
}

func readTodos() *[]Todo {
	var todos []Todo

	data, err := os.ReadFile("data/data.json")
	if err != nil {
		log.Println("Database file not found, creating")

		initialTodos := []Todo{
			{Id: 1, Title: "Task 1", Done: false},
			{Id: 2, Title: "Task 2", Done: true},
			{Id: 3, Title: "Task 3", Done: true},
		}
		writeTodos(&initialTodos)
		return &initialTodos
	}

	err = json.Unmarshal(data, &todos)
	if err != nil {
		log.Fatal("Failed to read from database file", err)
	}
	return &todos
}

func writeTodos(todos *[]Todo) {
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Fatal("Failed to marshall todos", err)
	}
	err = os.WriteFile("data/data.json", bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}

func IndexGetHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/index.gohtml"))
	var pageModel = IndexPageModel{
		PageTitle: "My TODO list",
	}

	var _ = indexTemplate.Execute(w, pageModel)
}

func TodosGetHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/todos.gohtml"))

	var todos = readTodos()

	var pageModel = TodosDTO{
		Todos: *todos,
	}

	var _ = indexTemplate.Execute(w, pageModel)
}

func ApiTodosPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		log.Println("Empty id passed")
		w.WriteHeader(400)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Unable to convert %s to int\n", idParam)
		w.WriteHeader(400)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		w.WriteHeader(400)
		return
	}
	var done = r.Form.Get("done") == "on"

	var todos = readTodos()

	for i := range *todos {
		if (*todos)[i].Id == id {
			(*todos)[i].Done = done
			break
		}
	}

	writeTodos(todos)

	w.WriteHeader(200)
}

func ApiTodosPutHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		w.WriteHeader(400)
		return
	}

	var title = r.Form.Get("title")
	if title == "" {
		w.WriteHeader(400)
		return
	}

	var todos = readTodos()
	var maxId = 1
	for _, todo := range *todos {
		if maxId < todo.Id {
			maxId = todo.Id
		}
	}
	maxId += 1

	*todos = append(*todos, Todo{Id: maxId, Title: title})

	writeTodos(todos)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(200)
}

func ApiTodosDelHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	if idParam == "" {
		log.Println("Empty id passed")
		w.WriteHeader(400)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Unable to convert %s to int\n", idParam)
		w.WriteHeader(400)
		return
	}

	var todos = readTodos()
	index := 0
	for i, todo := range *todos {
		if todo.Id == id {
			index = i
			break
		}
	}
	if index == 0 {
		log.Printf("Todo with id %d is not found", id)
		w.WriteHeader(400)
		return
	}

	*todos = removeAt(*todos, index)

	writeTodos(todos)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(200)
}

func removeAt[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
