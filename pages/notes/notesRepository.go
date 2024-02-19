package notes

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const dbFile = "data/notes.json"

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

func readNotes() *[]Note {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		log.Println("Database file not found, creating")
		initialNotes := make([]Note, 0)
		writeNotes(&initialNotes)
		return &initialNotes
	}

	var notes []Note
	err = json.Unmarshal(data, &notes)
	if err != nil {
		log.Fatal("Failed to read from database file", err)
	}
	return &notes
}

func writeNotes(notes *[]Note) {
	bytes, err := json.Marshal(notes)
	if err != nil {
		log.Fatal("Failed to marshall notes", err)
	}
	err = os.WriteFile(dbFile, bytes, 0666)
	if err != nil {
		log.Fatal("Failed to write database file", err)
	}
}
