package notes

import (
	"fmt"
	"go-server/helpers"
	"time"
)

func deleteNoteById(id int) error {
	notes := readNotes()

	index, note := helpers.FindByID(notes, id)
	if note == nil {
		return fmt.Errorf("note with id %d is not found", id)
	}

	*notes = helpers.RemoveAt(*notes, index)
	writeNotes(notes)

	return nil
}

// todo: update dto?
func updateNote(id int, title string, description string) error {
	notes := readNotes()
	_, note := helpers.FindByID(notes, id)
	if note == nil {
		return fmt.Errorf("note with id %d is not found", id)
	}

	note.Title = title
	note.Description = description
	note.Updated = time.Now()

	writeNotes(notes)

	return nil
}

func addNote(title string, description string) {
	notes := readNotes()

	*notes = append(*notes, Note{
		Id:          helpers.GenerateNextId(notes),
		Title:       title,
		Description: description,
		Created:     time.Now(),
		Updated:     time.Now(),
	})

	writeNotes(notes)
}
