package todos

import (
	"encoding/json"
	"log"
	"os"
)

const dbFile = "data/todos.json"

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t Todo) ID() int {
	return t.Id
}

func readTodos() *[]Todo {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		log.Println("Database file not found, creating")
		initialTodos := make([]Todo, 0)
		writeTodos(&initialTodos)
		return &initialTodos
	}

	var todos []Todo
	err = json.Unmarshal(data, &todos)
	if err != nil {
		log.Fatal("Failed to read from database file ", err)
	}
	return &todos
}

func writeTodos(todos *[]Todo) {
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Fatal("Failed to marshall todos", err)
	}
	err = os.WriteFile(dbFile, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
