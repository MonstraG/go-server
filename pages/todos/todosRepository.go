package todos

import (
	"encoding/json"
	"go-server/setup"
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

const dbFilePath = "todos.json"

type Repository struct {
	dbFilePath string
}

func NewRepository(config setup.AppConfig) Repository {
	return Repository{
		dbFilePath: config.DatabaseFolder + dbFilePath,
	}
}

func (repository Repository) readTodos() *[]Todo {
	data, err := os.ReadFile(repository.dbFilePath)
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
	err = os.WriteFile(repository.dbFilePath, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
