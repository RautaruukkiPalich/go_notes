package server

import (
	"encoding/json"
	"net/http"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

func (s *Server) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := s.getUserIdFromContext(r)

		limit, offset, filter_author, filter_body := getFiltersFromRequest(r)

		notes, err := s.store.Note().GetNotes(userID, filter_body, filter_author, limit, offset)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, errorResponse{http.StatusInternalServerError,	ErrInternalServerError.Error()})
			return
		}
		s.respond(w, r, http.StatusOK, notes)
	}
}

func (s *Server) PostNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := s.getUserIdFromContext(r)
		if userID == 0 {
			s.error(w, r, errorResponse{http.StatusForbidden, ErrNoPermissons.Error()})
			return
		}

		defer r.Body.Close()

		note := &model.Note{}
		if err := json.NewDecoder(r.Body).Decode(note); err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusUnprocessableEntity, ErrUnprocessableEntity.Error()})
			return
		}
		note.AuthorID = userID

		if err := s.store.Note().Set(note); err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusInternalServerError, ErrInternalServerError.Error()})
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *Server) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := getIdVarFromRequest(r)

		if err != nil {
			s.logger.Info(err.Error())
			s.error(w, r, errorResponse{http.StatusBadRequest, ErrParseURL.Error()})
			return
		}

		userID := s.getUserIdFromContext(r)
		if userID == 0 {
			s.error(w, r, errorResponse{http.StatusForbidden, ErrNoPermissons.Error()})
			return
		} 

		note, err := s.GetNoteById(id)
		if err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusNotFound, ErrNotFound.Error()})
			return
		}

		if !note.IsPublic && note.AuthorID != userID {
			s.error(w, r, errorResponse{http.StatusForbidden, ErrNoPermissons.Error()})
			return
		} 

		s.respond(w, r, http.StatusOK, note)
	}
}

func (s *Server) PatchNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdVarFromRequest(r)

		if err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusBadRequest, ErrParseURL.Error()})
			return
		}

		userID := s.getUserIdFromContext(r)

		note, err := s.store.Note().GetNoteById(id)
		if err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusNotFound, ErrNotFound.Error()})
			return
		}

		if note.AuthorID != userID {
			s.logger.Errorf("%d: forbidden patch note id%d", userID, note.ID)
			s.error(w, r, errorResponse{http.StatusForbidden, ErrNoPermissons.Error()})
			return
		}

		defer r.Body.Close()

		patchNote := &model.Note{}
		if err := json.NewDecoder(r.Body).Decode(patchNote); err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusUnprocessableEntity, ErrUnprocessableEntity.Error()})
			return
		}

		note.Body = patchNote.Body
		note.IsPublic = patchNote.IsPublic

		if err := s.store.Note().Patch(&note); err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusInternalServerError,	ErrInternalServerError.Error()})
			return
		}

		if err := s.cache.Note().Set(&note); err != nil {
			s.logger.Error(err.Error())
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *Server) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdVarFromRequest(r)

		if err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusBadRequest, ErrParseURL.Error()})
			return
		}

		userID := s.getUserIdFromContext(r)

		note, err := s.store.Note().GetNoteById(id)
		if err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusNotFound, ErrNotFound.Error()})
			return
		}

		if note.AuthorID != userID {
			s.error(w, r, errorResponse{http.StatusForbidden, ErrNoPermissons.Error()})
			return
		}

		if err = s.store.Note().Delete(id); err != nil {
			s.logger.Error(err.Error())
			s.error(w, r, errorResponse{http.StatusInternalServerError,	ErrInternalServerError.Error()})
			return
		}

		if err := s.cache.Note().Delete(id); err != nil {
			s.logger.Error(err.Error())
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
