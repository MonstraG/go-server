package helpers

import (
	"log"
	"net/http"
	"strconv"
)

// ParseFormRespondErr runs http.Request ParseForm, logs errors and writes http.StatusBadRequest to http.ResponseWriter
func ParseFormRespondErr(w MyWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		w.WriteHeader(http.StatusBadRequest)
	}
	return err
}

// ParseIdPathValueRespondErr tries to get 'id' from r.PathValue(),
// logging errors and writing 400 to http.ResponseWriter
func ParseIdPathValueRespondErr(w MyWriter, r *http.Request) int {
	idParam := r.PathValue("id")
	if idParam == "" {
		log.Println("Empty id passed")
		w.WriteHeader(http.StatusBadRequest)
		return 0
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Unable to convert %s to int\n", idParam)
		w.WriteHeader(http.StatusBadRequest)
		return 0
	}
	return id
}
