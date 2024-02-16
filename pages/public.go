package pages

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func PublicFolderHandler(w http.ResponseWriter, r *http.Request) {
	var pathValue = r.PathValue("path")
	path := fmt.Sprintf("public/%s", pathValue)
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println("Failed to read file", err)
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println("Failed to write to output???", err)
	}
}
