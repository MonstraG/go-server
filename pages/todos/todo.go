package todos

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t Todo) ID() int {
	return t.Id
}
