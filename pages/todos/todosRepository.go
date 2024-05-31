package todos

import (
	"go-server/helpers"
	"go-server/setup"
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
	return helpers.ReadData[Todo](repository.dbFilePath)
}

func (repository Repository) writeTodos(todos *[]Todo) {
	helpers.WriteData(repository.dbFilePath, todos)
}
