package todos

import (
	"fmt"
	"go-server/helpers"
)

func deleteTodoById(id int) error {
	todos := readTodos()

	index, todo := helpers.FindTodoByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	*todos = helpers.RemoveAt(*todos, index)
	writeTodos(todos)

	return nil
}

func addTodo(title string) {
	todos := readTodos()

	*todos = append(*todos, Todo{
		Id:    helpers.GenerateNextId(todos),
		Title: title,
	})

	writeTodos(todos)
}

func setTodoDoneState(id int, done bool) error {
	todos := readTodos()

	_, todo := helpers.FindTodoByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todo.Done = done
	writeTodos(todos)
	return nil
}
