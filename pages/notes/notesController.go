package notes

import (
	"encoding/json"
	"go-server/helpers"
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

var notesTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notes/notes.gohtml"))
var notesPageData = pages.PageData{
	PageTitle: "My note list",
}

type ListDTO struct {
	Notes []Note
}

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := notesTemplate.Execute(w, notesPageData)
	if err != nil {
		log.Fatal("Failed to render notes template", err)
	}
}

var notesListTemplate = template.Must(template.ParseFiles("pages/notes/notesList.gohtml"))

func GetListHandler(w http.ResponseWriter, _ *http.Request) {
	notes := readNotes()

	pageModel := ListDTO{
		Notes: *notes,
	}

	err := notesListTemplate.Execute(w, pageModel)
	if err != nil {
		log.Fatal("Failed to render notes template", err)
	}
}

var noteTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notes/note.gohtml"))

type notePageDTO struct {
	pages.PageData
	Note Note
}

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	notes := readNotes()
	_, note := helpers.FindByID(notes, id)

	data := notePageDTO{
		PageData: pages.PageData{PageTitle: note.Title},
		Note:     *note,
	}

	err := noteTemplate.Execute(w, data)
	if err != nil {
		log.Fatal("Failed to render notes template", err)
	}
}

func ApiGetHandler(w http.ResponseWriter, _ *http.Request) {
	notes := readNotes()
	bytes, err := json.Marshal(notes)
	if err != nil {
		log.Print("Failed to marshal notes", err)
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Print("Failed to write notes", err)
	}
}

// todo: actually like add update form and use this endpoint
func ApiPutHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	description := r.Form.Get("description")

	err = updateNote(id, title, description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ApiPostHandler(w http.ResponseWriter, r *http.Request) {
	err := helpers.ParseFormRespondErr(w, r)
	if err != nil {
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	description := r.Form.Get("description")

	addNote(title, description)

	w.Header().Set("HX-Trigger", "revalidateNotes")
	w.WriteHeader(http.StatusOK)
}

func ApiDelHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	err := deleteNoteById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
