package todos

import (
	"fmt"
	"go-server/helpers"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service Service) readTodos() *[]Todo {
	return service.repository.readTodos()
}

func (service Service) deleteTodoById(id int) error {
	todos := service.repository.readTodos()

	index, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	*todos = helpers.RemoveAt(*todos, index)
	service.repository.writeTodos(todos)

	return nil
}

func (service Service) addTodo(title string) {
	todos := service.repository.readTodos()

	*todos = append(*todos, Todo{
		Id:    helpers.GenerateNextId(todos),
		Title: title,
	})

	service.repository.writeTodos(todos)
}

func (service Service) setTodoDoneState(id int, done bool) error {
	todos := service.repository.readTodos()

	_, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todo.Done = done
	service.repository.writeTodos(todos)
	return nil
}
