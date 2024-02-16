package main

import (
	"example-go-server/pages"
	"fmt"
	"log"
	"net/http"
)

var address = ":8080"

func main() {
	http.HandleFunc("GET /{$}", pages.IndexHandler)

	log.Println(fmt.Sprintf("Starting server on http://localhost%s", address))
	log.Fatal(http.ListenAndServe(address, nil))
}
