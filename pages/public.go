package pages

import (
	"go-server/helpers"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// todo: add e-tags
// https://github.com/RHEnVision/provisioning-backend/blob/381501251cf9c09bdd2f860d43b113af271a792c/internal/middleware/etag_middleware.go

// todo: add gzip

func PublicHandler(w helpers.MyWriter, r *http.Request) {
	lw := helpers.MyWriter{ResponseWriter: w}
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Printf("Failed to stat file %s: %v", filename, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if fileInfo.IsDir() {
		log.Printf("Failed to stat file: %s, it's a directory", filename)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filename, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeContent(lw, r, filename, fileInfo.ModTime(), file)
}
