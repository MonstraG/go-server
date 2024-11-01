package todos

import (
	"fmt"
	"go-server/helpers"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service Service) readTodos() []Todo {
	return service.repository.readTodos()
}

func (service Service) deleteTodoById(id int) error {
	todos := service.repository.readTodos()

	index, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todos = helpers.RemoveAt(todos, index)
	service.repository.writeTodos(todos)

	return nil
}

func (service Service) addTodo(todo Todo) {
	todos := service.repository.readTodos()

	todo.Id = helpers.GenerateNextId(todos)

	todos = append(todos, todo)

	service.repository.writeTodos(todos)
}

func (service Service) setTodoDoneState(id int, done bool, updatedBy string) error {
	todos := service.repository.readTodos()

	_, todo := helpers.FindByID(todos, id)
	if todo == nil {
		return fmt.Errorf("todo with id %d is not found", id)
	}

	todo.Done = done
	todo.Updated = time.Now()
	todo.UpdatedBy = updatedBy
	service.repository.writeTodos(todos)
	return nil
}
