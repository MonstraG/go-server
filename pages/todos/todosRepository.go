package todos

import (
	"encoding/json"
	"go-server/helpers"
	"go-server/setup"
	"log"
)

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
	data := helpers.ReadData(repository.dbFilePath)

	var todos []Todo
	err := json.Unmarshal(data, &todos)
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

	helpers.WriteData(repository.dbFilePath, bytes)
}
