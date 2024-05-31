package notes

import (
	"go-server/helpers"
	"go-server/setup"
	"path"
)

const databaseFile = "notes.json"

type Repository struct {
	dbFilePath string
}

func NewRepository(config setup.AppConfig) Repository {
	return Repository{
		dbFilePath: path.Join(config.DatabaseFolder, databaseFile),
	}
}

func (repository Repository) readNotes() *[]Note {
	return helpers.ReadData[Note](repository.dbFilePath)
}

func (repository Repository) writeNotes(notes *[]Note) {
	helpers.WriteData(repository.dbFilePath, notes)
}
