package pages

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-server/helpers"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func PublicHandler(w helpers.MyWriter, r *helpers.MyRequest) {
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

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filename, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eTag := calculateETag(file)
	w.Header().Set("ETag", eTag)

	readSeeker := bytes.NewReader(file)

	// this handles ETag matches inside
	http.ServeContent(lw, &r.Request, filename, fileInfo.ModTime(), readSeeker)
}

var hasher = sha256.New()

// calculateETag generates a SHA-256 hash of the content and adds `W/` prefix to a hash to indicate weak comparison
func calculateETag(content []byte) string {
	hasher.Reset()
	hasher.Write(content)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("W/\"%s\"", hash)
}
