package notes

import (
	"encoding/json"
	"go-server/setup"
	"log"
	"os"
	"time"
)

type Note struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

func (n Note) ID() int {
	return n.Id
}

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
	data, err := os.ReadFile(repository.dbFilePath)
	if err != nil {
		log.Println("Database file not found, creating")
		initialNotes := make([]Note, 0)
		repository.writeNotes(&initialNotes)
		return &initialNotes
	}

	var notes []Note
	err = json.Unmarshal(data, &notes)
	if err != nil {
		log.Fatal("Failed to read from database file", err)
	}
	return &notes
}

func (repository Repository) writeNotes(notes *[]Note) {
	bytes, err := json.Marshal(notes)
	if err != nil {
		log.Fatal("Failed to marshall notes", err)
	}
	err = os.WriteFile(repository.dbFilePath, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
