package setup

import (
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

// CreateBasicAuthMiddleware returns middleware that requires basic auth
func CreateBasicAuthMiddleware(config AppConfig) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || config.Auth.Username != username || config.Auth.Password != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="server", charset="UTF-8"`)
				w.WriteHeader(401)
				return
			}

			next(w, r)
		}
	}
}
