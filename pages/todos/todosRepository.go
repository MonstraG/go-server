package todos

import (
	"encoding/json"
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

type DTO struct {
	Todos []Todo
}

// SetStateAction is an alias for function that changes todos and reports failures
type SetStateAction func(*[]Todo) bool

func deleteTodoAtIdAction(w http.ResponseWriter, id int) SetStateAction {
	return func(todos *[]Todo) bool {
		index, todo := findTodoById(todos, id)
		if todo == nil {
			log.Printf("Todo with id %d is not found", id)
			w.WriteHeader(400)
			return false
		}
		*todos = removeAt(*todos, index)
		return true
	}
}

// updateTodos reads from db, applies change and commits to db if successful
func updateTodos(change SetStateAction) {
	var todos = readTodos()

	ok := change(todos)
	if ok {
		writeTodos(todos)
	}
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
