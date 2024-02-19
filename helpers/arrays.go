package helpers

type Identifiable interface {
	ID() int
}

// RemoveAt removes element from slice at index, keeping order
func RemoveAt[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

// FindTodoByID returns index and element in slice together with a pointer to an element, allowing modifications
func FindTodoByID[T Identifiable](items *[]T, id int) (int, *T) {
	for i := range *items {
		if (*items)[i].ID() == id {
			return i, &(*items)[i]
		}
	}
	return 0, nil
}

// GenerateNextId finds next unoccupied id
func GenerateNextId[T Identifiable](items *[]T) int {
	maxId := 0
	for _, todo := range *items {
		if maxId < todo.ID() {
			maxId = todo.ID()
		}
	}
	return maxId + 1
}
