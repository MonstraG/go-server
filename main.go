package main

import (
	"go-server/pages"
	"log"
)

var address = ":8080"

func main() {
	app := NewApp()

	app.Use(LoggingMiddleware)

	app.HandleFunc("GET /{$}", pages.IndexGetHandler)
	app.HandleFunc("POST /api/todos/{id}", pages.ApiTodosPostHandler)
	app.HandleFunc("GET /todos", pages.TodosGetHandler)
	app.HandleFunc("PUT /api/todos", pages.ApiTodosPutHandler)
	app.HandleFunc("GET /public/{path...}", pages.PublicFolderHandler)

	log.Fatal(app.ListenAndServe(address))
}
