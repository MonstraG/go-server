package notes

import (
	"fmt"
	"go-server/helpers"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service Service) readNotes() []Note {
	return service.repository.readNotes()
}

func (service Service) deleteNoteById(id int) error {
	notes := service.repository.readNotes()

	index, note := helpers.FindByID(notes, id)
	if note == nil {
		return fmt.Errorf("note with id %d is not found", id)
	}

	notes = helpers.RemoveAt(notes, index)
	service.repository.writeNotes(notes)

	return nil
}

func (service Service) updateNote(id int, title string, description string, updatedBy string) error {
	notes := service.repository.readNotes()
	_, note := helpers.FindByID(notes, id)
	if note == nil {
		return fmt.Errorf("note with id %d is not found", id)
	}

	note.Title = title
	note.Description = description
	note.Updated = time.Now()
	note.UpdatedBy = updatedBy

	service.repository.writeNotes(notes)

	return nil
}

func (service Service) addNote(note Note) {
	notes := service.repository.readNotes()

	note.Id = helpers.GenerateNextId(notes)

	notes = append(notes, note)

	service.repository.writeNotes(notes)
}
