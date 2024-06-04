package helpers

import (
	"log"
	"net/http"
)

// MyWriter taken from https://stackoverflow.com/a/43976633/11593686
type MyWriter struct {
	http.ResponseWriter
}

// WriteSilent calls w.Write without telling you the result
func (w MyWriter) WriteSilent(p []byte) {
	_, err := w.ResponseWriter.Write(p)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func (w MyWriter) WriteStringSilent(p string) {
	w.WriteSilent([]byte(p))
}
