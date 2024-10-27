package pages

import (
	"go-server/helpers"
	"net/http"
	"path/filepath"
)

// todo: add e-tags
// https://github.com/RHEnVision/provisioning-backend/blob/381501251cf9c09bdd2f860d43b113af271a792c/internal/middleware/etag_middleware.go

// todo: add gzip

func PublicHandler(w helpers.MyWriter, r *http.Request) {
	lw := helpers.MyWriter{ResponseWriter: w}
	pathQueryParam := r.PathValue("path")
	file := filepath.Join("public", pathQueryParam)

	http.ServeFile(lw, r, file)
}
