package pages

import (
	"fmt"
	"go-server/helpers"
	"log"
	"net/http"
	"os"
	"strings"
)

// todo: add e-tags
// https://github.com/RHEnVision/provisioning-backend/blob/381501251cf9c09bdd2f860d43b113af271a792c/internal/middleware/etag_middleware.go

func PublicHandler(w helpers.MyWriter, r *http.Request) {
	lw := helpers.MyWriter{ResponseWriter: w}
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

	lw.WriteSilent(content)
}
