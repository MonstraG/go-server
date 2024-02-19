package main

import (
	"fmt"
	"log"
	"net/http"
)

// App = http.ServeMux + Middleware slice
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
	config      Config
}

// NewApp is a constructor for App
func NewApp(config Config) *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
		config:      config,
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
