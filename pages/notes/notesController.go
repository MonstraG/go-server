package notes

import (
	"go-server/helpers"
	"go-server/pages"
	"go-server/pages/notFound"
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

func (controller Controller) GetHandler(w *helpers.MyWriter, _ *helpers.MyRequest) {
	w.ExecuteTemplate(notesTemplate, notesPageData)
}

var notesListTemplate = template.Must(template.ParseFiles("pages/notes/tmpl/notesList.gohtml"))

type ListDTO struct {
	Notes []Note
}

func (controller Controller) GetListHandler(w *helpers.MyWriter, _ *helpers.MyRequest) {
	data := ListDTO{
		Notes: controller.service.readNotes(),
	}

	w.ExecuteTemplate(notesListTemplate, data)
}

var noteTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notes/tmpl/note.gohtml"))

type notePageDTO struct {
	pages.PageData
	Note Note
}

func (controller Controller) GetNoteHandler(w *helpers.MyWriter, r *helpers.MyRequest) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	notes := controller.service.readNotes()
	_, note := helpers.FindByID(notes, id)

	if note == nil {
		notFound.GetHandler(w, r)
		return
	}

	data := notePageDTO{
		PageData: pages.PageData{PageTitle: note.Title},
		Note:     *note,
	}

	w.ExecuteTemplate(noteTemplate, data)
}

var noteEditTemplate = template.Must(template.ParseFiles("pages/notes/tmpl/noteEdit.gohtml"))

func (controller Controller) GetNoteEditHandler(w *helpers.MyWriter, r *helpers.MyRequest) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}

	notes := controller.service.readNotes()
	_, note := helpers.FindByID(notes, id)

	if note == nil {
		notFound.GetHandler(w, r)
		return
	}

	w.ExecuteTemplate(noteEditTemplate, note)
}

func (controller Controller) ApiPutHandler(w *helpers.MyWriter, r *helpers.MyRequest) {
	id := helpers.ParseIdPathValueRespondErr(w, r)
	if id == 0 {
		return
	}
	ok := helpers.ParseFormRespondErr(w, r)
	if !ok {
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	description := r.Form.Get("description")

	err := controller.service.updateNote(id, title, description, r.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	w.Header().Set("HX-Redirect", "/notes")
	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiPostHandler(w *helpers.MyWriter, r *helpers.MyRequest) {
	ok := helpers.ParseFormRespondErr(w, r)
	if !ok {
		return
	}

	note := NewNote(r)
	if note == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	controller.service.addNote(*note)

	w.Header().Set("HX-Trigger", "revalidateNotes")
	w.WriteHeader(http.StatusOK)
}

func (controller Controller) ApiDelHandler(w *helpers.MyWriter, r *helpers.MyRequest) {
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

	w.Header().Set("HX-Redirect", "/notes")
	w.WriteHeader(http.StatusOK)
}
