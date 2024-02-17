package main

import (
	"fmt"
	"go-server/pages/notFound"
	"log"
	"net/http"
	"time"
)

// HandlerFn is an alias for http.HandlerFunc argument
type HandlerFn func(w http.ResponseWriter, r *http.Request)

// Middleware is just a HandlerFn that returns a HandlerFn
type Middleware func(HandlerFn) HandlerFn

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next HandlerFn) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// HtmxPartialMiddleware guards against direct browser navigations to partials
// It returns notFound if request wasn't made by htmx (Hx-Request header)
func HtmxPartialMiddleware(next HandlerFn) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		isHtmxRequest := r.Header.Get("Hx-Request") == "true"
		if !isHtmxRequest {
			notFound.RedirectToNotFoundHandler(w, r)
			return
		}

		next(w, r)
	}
}

// App = http.ServeMux + Middleware slice
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewApp is a constructor for App
func NewApp() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
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
func (app *App) ListenAndServe(address string) error {
	log.Println(fmt.Sprintf("Starting server on http://localhost%s", address))
	return http.ListenAndServe(address, app.mux)
}
