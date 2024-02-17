package todos

import (
	"encoding/json"
	"log"
	"os"
)

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// SetStateAction is an alias for function that changes todos and reports failures
type SetStateAction func(*[]Todo) bool

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
	maxId := 0
	for _, todo := range *todos {
		if maxId < todo.Id {
			maxId = todo.Id
		}
	}
	return maxId + 1
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
