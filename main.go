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

	// todo: can we make this automatic?
	// the Injection part is manual here)
	todosController := todos.NewController(todos.NewService(todos.NewRepository(config)))
	notesController := notes.NewController(notes.NewService(notes.NewRepository(config)))

	// I'm torn between registering all endpoints in a single file vs in each controller separately.
	// It seems like doing this in a single file allows to somewhat more easily spot conflicts or inconsistencies.
	// Probably for a bigger project, this will have to be done per-controller.

	// pages
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /404", notFound.GetHandler)
	app.HandleFunc("GET /*", notFound.RedirectToNotFoundHandler)
	app.HandleFunc("GET /todos", todosController.GetHandler)
	app.HandleFunc("GET /notes", notesController.GetHandler)
	app.HandleFunc("GET /notes/{id}", notesController.GetNoteHandler)

	// partials
	app.HandleFunc("GET /todosList", setup.HtmxPartialMiddleware(todosController.GetListHandler))
	app.HandleFunc("GET /notesList", setup.HtmxPartialMiddleware(notesController.GetListHandler))
	app.HandleFunc("GET /notes/{id}/edit", setup.HtmxPartialMiddleware(notesController.GetNoteEditHandler))

	// api
	app.HandleFunc("PUT /api/todos/{id}", todosController.ApiPutHandler)
	app.HandleFunc("POST /api/todos", todosController.ApiPostHandler)
	app.HandleFunc("DELETE /api/todos/{id}", todosController.ApiDelHandler)

	app.HandleFunc("PUT /api/notes/{id}", notesController.ApiPutHandler)
	app.HandleFunc("POST /api/notes", notesController.ApiPostHandler)
	app.HandleFunc("DELETE /api/notes/{id}", notesController.ApiDelHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	log.Fatal(app.ListenAndServe())
}
