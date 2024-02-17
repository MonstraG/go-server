package main

import (
	"go-server/pages"
	"log"
)

var address = ":8080"

func main() {
	app := NewApp()

	app.Use(LoggingMiddleware)

	// pages or partials
	app.HandleFunc("GET /{$}", pages.IndexGetHandler)
	app.HandleFunc("GET /todos", pages.TodosGetHandler)

	// api
	app.HandleFunc("POST /api/todos/{id}", pages.ApiTodosPostHandler)
	app.HandleFunc("PUT /api/todos", pages.ApiTodosPutHandler)
	app.HandleFunc("DELETE /api/todos/{id}", pages.ApiTodosDelHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicFolderHandler)

	log.Fatal(app.ListenAndServe(address))
}
