package setup

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// App = http.ServeMux + Middleware slice
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
	config      AppConfig
}

// NewApp is a constructor for App
func NewApp(appConfig AppConfig) *App {
	err := os.MkdirAll(appConfig.DatabaseFolder, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create database folder")
	}

	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
		config:      appConfig,
	}
}

// Use adds Middleware to chain
func (app *App) Use(mw Middleware) {
	app.middlewares = append(app.middlewares, mw)
}

// HandleFunc is a wrapper around normal http.HandleFunc but calling all Middleware-s first
func (app *App) HandleFunc(pattern string, handlerFunc HandlerFn) {
	app.mux.HandleFunc(pattern, applyMiddlewares(handlerFunc, app.middlewares))
}

// applyMiddlewares runs all middlewares in order
func applyMiddlewares(h HandlerFn, middlewares []Middleware) HandlerFn {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// ListenAndServe is a wrapper around normal http.ListenAndServe
func (app *App) ListenAndServe() error {
	log.Println(fmt.Sprintf("Starting server on %s", app.config.Host))
	return http.ListenAndServe(app.config.Host, app.mux)
}
