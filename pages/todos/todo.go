package todos

import (
	"go-server/helpers"
	"time"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	Created   time.Time
	Updated   time.Time
	UpdatedBy string
}

func (t Todo) ID() int {
	return t.Id
}

func NewTodo(r *helpers.MyRequest) *Todo {
	title := r.Form.Get("title")
	if title == "" {
		return nil
	}

	return &Todo{
		Title:     title,
		Created:   time.Now(),
		Updated:   time.Now(),
		UpdatedBy: r.Username,
	}
}
