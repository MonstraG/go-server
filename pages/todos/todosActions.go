package todos

import (
	"fmt"
)

func deleteTodoById(id int) error {
	todos := readTodos()

	index, todo := findTodoById(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	*todos = removeAt(*todos, index)
	writeTodos(todos)

	return nil
}

func addTodo(title string) {
	todos := readTodos()

	*todos = append(*todos, Todo{
		Id:    generateNextId(todos),
		Title: title,
	})

	writeTodos(todos)
}

func setTodoDoneState(id int, done bool) error {
	todos := readTodos()

	_, todo := findTodoById(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todo.Done = done
	writeTodos(todos)
	return nil
}
