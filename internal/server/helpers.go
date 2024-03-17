package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type userKey string
const (
	UserKey userKey = "userID" 
)

func getIdVarFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	strID := vars["id"]
	return strconv.Atoi(strID)
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