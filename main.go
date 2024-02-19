package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/notFound"
	"go-server/pages/notes"
	"go-server/pages/todos"
	"go-server/setup"
	"log"
)

func main() {
	config := setup.ReadConfig()

	app := setup.NewApp(config)

	app.Use(setup.LoggingMiddleware)

	todosController := todos.NewController(todos.NewService(todos.NewRepository(config)))

	// pages
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /404", notFound.GetHandler)
	app.HandleFunc("GET /*", notFound.RedirectToNotFoundHandler)
	app.HandleFunc("GET /todos", todosController.GetHandler)
	app.HandleFunc("GET /notes", notes.GetHandler)
	app.HandleFunc("GET /notes/{id}", notes.GetNoteHandler)

	// partials
	app.HandleFunc("GET /todosList", setup.HtmxPartialMiddleware(todosController.GetListHandler))
	app.HandleFunc("GET /notesList", setup.HtmxPartialMiddleware(notes.GetListHandler))

	// api
	app.HandleFunc("GET /api/todos", todosController.ApiGetHandler)
	app.HandleFunc("PUT /api/todos/{id}", todosController.ApiPutHandler)
	app.HandleFunc("POST /api/todos", todosController.ApiPostHandler)
	app.HandleFunc("DELETE /api/todos/{id}", todosController.ApiDelHandler)

	app.HandleFunc("GET /api/notes", notes.ApiGetHandler)
	app.HandleFunc("PUT /api/notes/{id}", notes.ApiPutHandler)
	app.HandleFunc("POST /api/notes", notes.ApiPostHandler)
	app.HandleFunc("DELETE /api/notes/{id}", notes.ApiDelHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	log.Fatal(app.ListenAndServe())
}
