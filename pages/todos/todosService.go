package todos

import (
	"fmt"
	"go-server/helpers"
)

type Service struct {
	Repository Repository
}

func (service Service) deleteTodoById(id int) error {
	todos := service.Repository.readTodos()

	index, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	*todos = helpers.RemoveAt(*todos, index)
	service.Repository.writeTodos(todos)

	return nil
}

func (service Service) addTodo(title string) {
	todos := service.Repository.readTodos()

	*todos = append(*todos, Todo{
		Id:    helpers.GenerateNextId(todos),
		Title: title,
	})

	service.Repository.writeTodos(todos)
}

func (service Service) setTodoDoneState(id int, done bool) error {
	todos := service.Repository.readTodos()

	_, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todo.Done = done
	service.Repository.writeTodos(todos)
	return nil
}
