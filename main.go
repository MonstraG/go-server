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

	// Add routes
	app.HandleFunc("GET /{$}", pages.IndexHandler)
	app.HandleFunc("GET /public/{path...}", pages.PublicFolderHandler)

	log.Fatal(app.ListenAndServe(address))
}
