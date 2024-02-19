package notes

import (
	"encoding/json"
	"go-server/helpers"
	"go-server/pages"
	"html/template"
	"log"
	"net/http"
)

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

var notesTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notes/tmpl/notes.gohtml"))
var notesPageData = pages.PageData{
	PageTitle: "My note list",
}

type ListDTO struct {
	Notes []Note
}

func (controller Controller) GetHandler(w http.ResponseWriter, _ *http.Request) {
	err := notesTemplate.Execute(w, notesPageData)
	if err != nil {
		log.Fatal("Failed to render notes template", err)
	}
}

var notesListTemplate = template.Must(template.ParseFiles("pages/notes/tmpl/notesList.gohtml"))

func (controller Controller) GetListHandler(w http.ResponseWriter, _ *http.Request) {
	notes := controller.service.readNotes()

	pageModel := ListDTO{
		Notes: *notes,
	}

	err := notesListTemplate.Execute(w, pageModel)
	if err != nil {
		log.Fatal("Failed to render notes template", err)
	}
}

var noteTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notes/tmpl/note.gohtml"))

type notePageDTO struct {
	pages.PageData
	Note Note
}

func (controller Controller) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	notes := controller.service.readNotes()
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

func (controller Controller) ApiGetHandler(w http.ResponseWriter, _ *http.Request) {
	notes := controller.service.readNotes()
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
func (controller Controller) ApiPutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = controller.service.updateNote(id, title, description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiPostHandler(w http.ResponseWriter, r *http.Request) {
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

	controller.service.addNote(title, description)

	w.Header().Set("HX-Trigger", "revalidateNotes")
	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiDelHandler(w http.ResponseWriter, r *http.Request) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	err := controller.service.deleteNoteById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
