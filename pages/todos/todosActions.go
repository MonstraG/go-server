package todos

import (
	"log"
	"net/http"
)

// deleteTodoAtIdAction runs a SetStateAction to delete a todo by id
func deleteTodoAtIdAction(w http.ResponseWriter, id int) SetStateAction {
	return func(todos *[]Todo) bool {
		index, todo := findTodoById(todos, id)
		if todo == nil {
			log.Printf("Todo with id %d is not found", id)
			w.WriteHeader(http.StatusBadRequest)
			return false
		}
		*todos = removeAt(*todos, index)
		return true
	}
}

// addTodoAction runs a SetStateAction to add a todo
func addTodoAction(title string) SetStateAction {
	return func(todos *[]Todo) bool {
		*todos = append(*todos, Todo{
			Id:    generateNextId(todos),
			Title: title,
		})
		return true
	}
}

// changeDoneAction runs a SetStateAction to change Todo.Done status by id
func changeDoneAction(w http.ResponseWriter, id int, done bool) SetStateAction {
	return func(todos *[]Todo) bool {
		_, todo := findTodoById(todos, id)
		if todo == nil {
			log.Printf("Todo with id %d is not found", id)
			w.WriteHeader(http.StatusBadRequest)
			return false
		}
		todo.Done = done
		return true
	}
}

// runTodosAction reads from db, applies change and commits to db if successful
func runTodosAction(change SetStateAction) {
	var todos = readTodos()

	ok := change(todos)
	if ok {
		writeTodos(todos)
	}
}
