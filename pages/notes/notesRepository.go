package notes

import (
	"go-server/helpers"
	"go-server/setup"
)

const dbFilePath = "notes.json"

type Repository struct {
	dbFilePath string
}

func NewRepository(config setup.AppConfig) Repository {
	return Repository{
		dbFilePath: config.DatabaseFolder + dbFilePath,
	}
}

func (repository Repository) readNotes() *[]Note {
	return helpers.ReadData[Note](repository.dbFilePath)
}

func (repository Repository) writeNotes(notes *[]Note) {
	helpers.WriteData(repository.dbFilePath, notes)
}
