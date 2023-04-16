package routers

import (
	"flipr_assignment/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterNoteRoutes(router *mux.Router) {
	router.Handle("/api/note/", (http.HandlerFunc(controllers.GetNotes))).Methods("GET")
	router.Handle("/api/note/{id}", (http.HandlerFunc(controllers.GetNote))).Methods("GET")
	router.Handle("/api/note/add/",(http.HandlerFunc(controllers.AddNote))).Methods("POST")
	router.Handle("/api/note/update/{id}", (http.HandlerFunc(controllers.UpdateNote))).Methods("PUT")
	router.Handle("/api/note/delete/{id}", (http.HandlerFunc(controllers.DeleteNote))).Methods("DELETE")
}
