package query

import (
	"log"
	"net/http"
	"strconv"
)

// ParseFormRespondErr runs http.Request ParseForm, logging errors and writing 400 to http.ResponseWriter
func ParseFormRespondErr(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		w.WriteHeader(400)
	}
	return err
}

// ParseIdPathValueRespondErr tries to get 'id' from r.PathValue(),
// logging errors and writing 400 to http.ResponseWriter
func ParseIdPathValueRespondErr(w http.ResponseWriter, r *http.Request) int {
	idParam := r.PathValue("id")
	if idParam == "" {
		log.Println("Empty id passed")
		w.WriteHeader(400)
		return 0
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Unable to convert %s to int\n", idParam)
		w.WriteHeader(400)
		return 0
	}
	return id
}
