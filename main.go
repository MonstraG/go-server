package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/notFound"
	"go-server/pages/todos"
	"log"
)

const address = ":8080"

func main() {
	app := NewApp()

	app.Use(LoggingMiddleware)

	// pages
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /404", notFound.GetHandler)
	app.HandleFunc("GET /*", notFound.RedirectToNotFoundHandler)

	// partials
	app.HandleFunc("GET /todos", HtmxPartialMiddleware(todos.GetHandler))

	// api
	app.HandleFunc("GET /api/todos", todos.ApiGetHandler)
	app.HandleFunc("POST /api/todos/{id}", todos.ApiPostHandler)
	app.HandleFunc("PUT /api/todos", todos.ApiPutHandler)
	app.HandleFunc("DELETE /api/todos/{id}", todos.ApiDelHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	log.Fatal(app.ListenAndServe(address))
}
