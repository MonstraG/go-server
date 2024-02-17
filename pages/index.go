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

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
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

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("pages/index.html"))
	var todos = readTodos()
	var pageModel = TodoPageData{
		PageTitle: "My TODO list",
		Todos:     *todos,
	}

	var _ = indexTemplate.Execute(w, pageModel)
}

func ApiTodoPost(w http.ResponseWriter, r *http.Request) {
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

	for i := range *todos {
		if (*todos)[i].Id == id {
			(*todos)[i].Done = !(*todos)[i].Done
			break
		}
	}

	writeTodos(todos)

	w.WriteHeader(200)
}
