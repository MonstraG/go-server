package helpers

import (
	"compress/gzip"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
)

// MyWriter taken from https://stackoverflow.com/a/43976633/11593686
type MyWriter struct {
	http.ResponseWriter

	supportsGzipEncoding bool
	gzipWriter           *gzip.Writer
}

func NewMyWriter(w http.ResponseWriter, r *http.Request) MyWriter {
	return MyWriter{ResponseWriter: w, supportsGzipEncoding: supportsGzipEncoding(r)}
}

func supportsGzipEncoding(r *http.Request) bool {
	acceptedEncodingHeader := r.Header.Get("Accept-Encoding")
	acceptedEncodings := strings.Split(acceptedEncodingHeader, ",")
	return slices.Contains(acceptedEncodings, gzipEncoding)
}

const gzipEncoding = "gzip"

// WriteSilent calls w.Write without telling you the result
func (w *MyWriter) WriteSilent(content []byte) {
	_, err := w.Write(content)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func (w *MyWriter) needsGzip() bool {
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		log.Println("Did not set Content-Type before writing!")
	}

	if strings.HasPrefix(contentType, "image/") && contentType != "image/svg+xml" {
		return false
	}
	return true
}

func (w *MyWriter) Write(content []byte) (int, error) {
	if !w.needsGzip() || !w.supportsGzipEncoding {
		return w.ResponseWriter.Write(content)
	}

	if w.gzipWriter == nil {
		w.Header().Set("Content-Encoding", gzipEncoding)
		w.gzipWriter = gzip.NewWriter(w.ResponseWriter)
	}

	return w.gzipWriter.Write(content)
}

// Close should be called when you're done writing all the content
func (w *MyWriter) Close() error {
	if w.gzipWriter != nil {
		return w.gzipWriter.Close()
	}
	return nil
}

func (w *MyWriter) SetContentTypeHTML() {
	w.Header().Set("Content-Type", "text/html")
}

func (w *MyWriter) ExecuteTemplate(template *template.Template, data any) {
	w.SetContentTypeHTML()
	err := template.Execute(w, data)
	if err != nil {
		log.Println("Failed to render template:\n", err)
	}
}
