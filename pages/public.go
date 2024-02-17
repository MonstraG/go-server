package pages

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	var pathValue = r.PathValue("path")
	path := fmt.Sprintf("public/%s", pathValue)

	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Failed to read file", err)
	}

	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println("Failed to write file to response", err)
	}
}
