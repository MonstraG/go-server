package main

import (
	"example-go-server/pages"
	"log"
)

var address = ":8080"

func main() {
	app := NewApp()

	// Add middleware
	app.Use(LoggingMiddleware)
	app.Use(AnotherMiddleware)

	// Add routes
	app.HandleFunc("GET /{$}", pages.IndexHandler)

	log.Fatal(app.ListenAndServe(address))
}
