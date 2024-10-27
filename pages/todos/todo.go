package todos

import "time"

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
