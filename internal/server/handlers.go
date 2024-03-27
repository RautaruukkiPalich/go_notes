package server

import (
	"encoding/json"
	"net/http"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

//	@Summary		Get Notes
//	@Security		ApiKeyAuth
//	@Description	get notes
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			limit			query		int		false	"limit"
//	@Param			offset			query		int		false	"offset"
//	@Param			filter_author	query		int		false	"filter_author"
//	@Param			filter_body		query		string	false	"filter_body"
//	@Success		200				{object}	[]model.Note
//	@Failure		400,404			{object}	errorResponse
//	@Success		500				{object}	errorResponse
//	@Success		default			{object}	errorResponse
//	@Router			/notes [get]
func (s *Server) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := s.getUserIdFromContext(r)

		limit, offset, filter_author, filter_body := getFiltersFromRequest(r)

		notes, err := s.store.Note().GetNotes(userID, filter_body, filter_author, limit, offset)
		if err != nil {
			s.logger.Error(err)
			s.error(w, r, errorResponse{http.StatusInternalServerError, ErrInternalServerError.Error()})
			return
		}
		s.JSONrespond(w, r, http.StatusOK, notes)
	}
}

//	@Summary		Post New Note
//	@Security		ApiKeyAuth
//	@Description	post note
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			input	body		notePostForm	true	"note"
//	@Success		200		{object}	nil
//	@Failure		400,404	{object}	errorResponse
//	@Success		500		{object}	errorResponse
//	@Success		default	{object}	errorResponse
//	@Router			/notes [post]
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

		s.JSONrespond(w, r, http.StatusOK, nil)
	}
}

//	@Summary		Get note by ID
//	@Security		ApiKeyAuth
//	@Description	get note
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"note identifier"
//	@Success		200		{object}	model.Note
//	@Failure		400,404	{object}	errorResponse
//	@Success		500		{object}	errorResponse
//	@Success		default	{object}	errorResponse
//	@Router			/notes/{id} [get]
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

		s.JSONrespond(w, r, http.StatusOK, note)
	}
}

//	@Summary		Patch note by ID
//	@Security		ApiKeyAuth
//	@Description	patch note
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"note identifier"
//	@Param			input	body		notePostForm	true	"note"
//	@Success		200		{object}	nil
//	@Failure		400,404	{object}	errorResponse
//	@Success		500		{object}	errorResponse
//	@Success		default	{object}	errorResponse
//	@Router			/notes/{id} [patch]
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
			s.logger.Errorf("user %d: forbidden patch note id%d", userID, note.ID)
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
			s.error(w, r, errorResponse{http.StatusInternalServerError, ErrInternalServerError.Error()})
			return
		}

		if err := s.cache.Note().Set(&note); err != nil {
			s.logger.Error(err.Error())
		}

		s.JSONrespond(w, r, http.StatusOK, nil)
	}
}

//	@Summary		Del note by ID
//	@Security		ApiKeyAuth
//	@Description	del note
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"note identifier"
//	@Success		200		{object}	nil
//	@Failure		400,404	{object}	errorResponse
//	@Success		500		{object}	errorResponse
//	@Success		default	{object}	errorResponse
//	@Router			/notes/{id} [delete]
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
			s.error(w, r, errorResponse{http.StatusInternalServerError, ErrInternalServerError.Error()})
			return
		}

		if err := s.cache.Note().Delete(id); err != nil {
			s.logger.Error(err.Error())
		}

		s.JSONrespond(w, r, http.StatusOK, nil)
	}
}
