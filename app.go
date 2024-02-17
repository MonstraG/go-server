package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)

// LoggingMiddleware logs the request details
func LoggingMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// App struct to hold our routes and middleware
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewApp creates and returns a new App with an initialized ServeMux and middleware slice
func NewApp() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
	}
}

// Use adds middleware to the chain
func (app *App) Use(mw Middleware) {
	app.middlewares = append(app.middlewares, mw)
}

// HandleFunc registers a handler for a specific route, applying all middlewares
func (app *App) HandleFunc(pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	app.mux.HandleFunc(pattern, applyMiddleware(handlerFunc, app.middlewares))
}

// ApplyMiddleware applies multiple middleware to a http.Handler
func applyMiddleware(h func(w http.ResponseWriter, r *http.Request), middlewares []Middleware) func(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func (app *App) ListenAndServe(address string) error {
	log.Println(fmt.Sprintf("Starting server on http://localhost%s", address))
	return http.ListenAndServe(address, app.mux)
}
