package notes

import (
	"encoding/json"
	"go-server/helpers"
	"go-server/setup"
	"log"
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
	data := helpers.ReadData(repository.dbFilePath)

	var notes []Note
	err := json.Unmarshal(data, &notes)
	if err != nil {
		log.Fatal("Failed to read from database file:\n", err)
	}
	return &notes
}

func (repository Repository) writeNotes(notes *[]Note) {
	bytes, err := json.Marshal(notes)
	if err != nil {
		log.Fatal("Failed to marshall notes:\n", err)
	}

	helpers.WriteData(repository.dbFilePath, bytes)
}
