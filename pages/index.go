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

// readTodos reads todos from data.json file
func readTodos() *[]Todo {
	var todos []Todo

	data, err := os.ReadFile("data/data.json")
	if err != nil {
		log.Println("Database file not found, creating")

		initialTodos := make([]Todo, 0)
		writeTodos(&initialTodos)
		return &initialTodos
	}

	err = json.Unmarshal(data, &todos)
	if err != nil {
		log.Fatal("Failed to read from database file", err)
	}
	return &todos
}

// writeTodos writes todos to data.json file
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

func ApiTodosPutHandler(w http.ResponseWriter, r *http.Request) {
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

// parseIdPathValueSendErr tries to get 'id' from r.PathValue(), logging errors and writing 400 to http.ResponseWriter
func parseIdPathValueSendErr(w http.ResponseWriter, r *http.Request) int {
	idParam := r.PathValue("id")
	if idParam == "" {
		log.Println("Empty id passed")
		w.WriteHeader(400)
		return 0
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Unable to convert %s to int\n", idParam)
		w.WriteHeader(400)
		return 0
	}
	return id
}

// parseFormSendErr runs http.Request ParseForm, logging errors and writing 400 to http.ResponseWriter
func parseFormSendErr(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		w.WriteHeader(400)
	}
	return err
}

func ApiTodosDelHandler(w http.ResponseWriter, r *http.Request) {
	id := parseIdPathValueSendErr(w, r)
	if id == 0 {
		return
	}

	var todos = readTodos()

	index, todo := findTodoById(todos, id)
	if todo == nil {
		log.Printf("Todo with id %d is not found", id)
		w.WriteHeader(400)
		return
	}
	*todos = removeAt(*todos, index)

	writeTodos(todos)

	w.Header().Set("HX-Trigger", "revalidateTodos")
	w.WriteHeader(200)
}

// removeAt removes element from slice at index, keeping order
func removeAt[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// findTodoById returns index and element in slice together with a pointer to an element, allowing modifications
func findTodoById(todos *[]Todo, id int) (int, *Todo) {
	for i := range *todos {
		if (*todos)[i].Id == id {
			return i, &(*todos)[i]
		}
	}
	return 0, nil
}

// generateNextId finds next unoccupied id
func generateNextId(todos *[]Todo) int {
	var maxId = 0
	for _, todo := range *todos {
		if maxId < todo.Id {
			maxId = todo.Id
		}
	}
	return maxId + 1
}
