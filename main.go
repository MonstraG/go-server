package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/notFound"
	"go-server/pages/todos"
	"log"
)

func main() {
	config := ReadConfig()

	app := NewApp(config)

	app.Use(LoggingMiddleware)

	// pages
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /404", notFound.GetHandler)
	app.HandleFunc("GET /*", notFound.RedirectToNotFoundHandler)
	app.HandleFunc("GET /todos", todos.GetHandler)

	// partials
	app.HandleFunc("GET /todosList", HtmxPartialMiddleware(todos.GetListHandler))

	// api
	app.HandleFunc("GET /api/todos", todos.ApiGetHandler)
	app.HandleFunc("PUT /api/todos/{id}", todos.ApiPutHandler)
	app.HandleFunc("POST /api/todos", todos.ApiPostHandler)
	app.HandleFunc("DELETE /api/todos/{id}", todos.ApiDelHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	log.Fatal(app.ListenAndServe())
}
