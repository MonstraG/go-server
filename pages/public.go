package pages

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	pathQueryParam := r.PathValue("path")
	path := fmt.Sprintf("public/%s", pathQueryParam)

	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Failed to read file", err)
	}

	if strings.HasSuffix(path, ".ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	}
	if strings.HasSuffix(path, ".svg") {
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	}
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println("Failed to write file to response", err)
	}
}
