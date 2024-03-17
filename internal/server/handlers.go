package server

import (
	"fmt"
	"net/http"
)

func (s *Server) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Get Notes")
	}
}

func (s *Server) PostNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Post Note")
	}
}

func (s *Server) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := getIdVarFromRequest(r)

		
		if err != nil {
			s.logger.Info(err.Error())
			s.error(w, r, errorResponse{
				http.StatusBadRequest,
				ErrParseURL.Error(),
			})
			return
		}

		note, err := s.store.Note().GetNoteById(id)
		if err != nil {
			s.error(w, r, errorResponse{
				http.StatusNotFound,
				ErrNotFound.Error(),
			})
			return
		}
		userID := s.getUserIdFromContext(r)

		if note.AuthorID != userID || userID == 0 || note.AuthorID == 0 {
			s.error(w, r, errorResponse{
				http.StatusBadRequest,
				ErrNoPermissons.Error(),
			})
			return
		}

		s.respond(w, r, http.StatusOK, note)
	}
}

func (s *Server) PatchNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Patch Note")
	}
}

func (s *Server) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Delete Note")
	}
}
