package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rautaruukkipalich/go_notes/internal/model"
)

type (
	userKey string
	
	NotePostForm struct {
		Body string `json:"body"`
		IsPublic bool `json:"is_public"`
	}
)

const (
	UserKey userKey = "userID" 
)

func getIdVarFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	strID := vars["id"]
	return strconv.Atoi(strID)
}


//limit, offset and body/author filters
func getFiltersFromRequest(r *http.Request) (int, int, int, string) {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}

	filter_author, err := strconv.Atoi(r.FormValue("filter_author"))
	if err != nil {
		filter_author = 0
	}

	filterBody := r.FormValue("filter_body")

	return limit, offset, filter_author, filterBody
}

func (s *Server) GetNoteById(id int) (model.Note, error) {
	var note model.Note

	s.logger.Infof("get note %d from cache...", id)
	note, err := s.cache.Note().GetNoteById(id)
	if err == nil {
		return note, nil
	}

	s.logger.Infof("get note %d from db...", id)
	note, err = s.store.Note().GetNoteById(id)
	if err != nil {
		return note, err
	}

	s.logger.Infof("set note %d to cache...", id)
	if err := s.cache.Note().Set(&note); err != nil {
		return note, err
	}
	return note, err
}

func (s *Server) getUserIdFromContext(r *http.Request) int {
	ctx := r.Context()
	id, ok := ctx.Value(UserKey).(int)

	if !ok {
		return 0
	}
	return id
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, err errorResponse) {
	s.logger.Printf("error: %v", err)
	s.respond(w, r, err.Code, err)
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Printf("error: %v", err)
		}
	}
}