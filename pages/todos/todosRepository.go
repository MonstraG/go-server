package todos

import (
	"go-server/helpers"
	"go-server/setup"
	"path/filepath"
)

const databaseFile = "todos.json"

type Repository struct {
	dbFilePath string
}

func NewRepository(config setup.AppConfig) Repository {
	return Repository{
		dbFilePath: filepath.Join(config.DatabaseFolder, databaseFile),
	}
}

func (repository Repository) readTodos() []Todo {
	return helpers.ReadData[Todo](repository.dbFilePath)
}

func (repository Repository) writeTodos(todos []Todo) {
	helpers.WriteData(repository.dbFilePath, todos)
}
