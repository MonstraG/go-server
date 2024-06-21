package dictionary

type Entry struct {
	Id          int    `json:"id"`
	Original    string `json:"original"`
	Translation string `json:"translation"`
}

func (e Entry) ID() int {
	return e.Id
}
