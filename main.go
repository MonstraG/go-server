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

	var authMiddleware = setup.CreateBasicAuthMiddleware(config)
	app.Use(authMiddleware)

	app.Use(setup.LoggingMiddleware)

	// Pure DI
	// Injection can be made automatic, but it turns out, it's just worse, esp. at this scale
	todosController := todos.NewController(todos.NewService(todos.NewRepository(config)))
	notesController := notes.NewController(notes.NewService(notes.NewRepository(config)))

	// I'm torn between registering all endpoints in a single file vs in each controller separately.
	// It seems like doing this in a single file allows to somewhat more easily spot conflicts or inconsistencies.
	// Probably for a bigger project, this will have to be done per-controller.

	// pages
	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /todos", todosController.GetHandler)
	app.HandleFunc("GET /notes", notesController.GetHandler)
	app.HandleFunc("GET /notes/{id}", notesController.GetNoteHandler)

	// partials
	app.HandleFunc("GET /todosList", setup.HtmxPartialMiddleware(todosController.GetListHandler))
	app.HandleFunc("GET /notesList", setup.HtmxPartialMiddleware(notesController.GetListHandler))
	app.HandleFunc("GET /notes/{id}/edit", setup.HtmxPartialMiddleware(notesController.GetNoteEditHandler))

	// It feels like htmx wants me to get rid of /api/ endpoints, or, more precisely,
	// remove `/api` from the url and return partials in them, but I'm not decided on that at the moment

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
