package notes

import (
	"go-server/helpers"
	"time"
)

type Note struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	UpdatedBy   string
}

func (n Note) ID() int {
	return n.Id
}

func NewNote(r *helpers.MyRequest) *Note {
	title := r.Form.Get("title")
	if title == "" {
		return nil
	}

	description := r.Form.Get("description")

	return &Note{
		Title:       title,
		Description: description,
		Created:     time.Now(),
		Updated:     time.Now(),
		UpdatedBy:   r.Username,
	}
}
