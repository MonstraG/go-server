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

func (t Todo) ID() int {
	return t.Id
}

type Repository struct {
	DatabaseFolder string
}

const dbFilePath = "todos.json"

func (repository Repository) getFilePath() string {
	return repository.DatabaseFolder + dbFilePath
}

func (repository Repository) readTodos() *[]Todo {
	data, err := os.ReadFile(repository.getFilePath())
	if err != nil {
		log.Println("Database file not found, creating")
		initialTodos := make([]Todo, 0)
		repository.writeTodos(&initialTodos)
		return &initialTodos
	}

	var todos []Todo
	err = json.Unmarshal(data, &todos)
	if err != nil {
		log.Fatal("Failed to read from database file ", err)
	}
	return &todos
}

func (repository Repository) writeTodos(todos *[]Todo) {
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Fatal("Failed to marshall todos", err)
	}
	err = os.WriteFile(repository.getFilePath(), bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
