package server

import "net/http"

func (s *Server) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}


func (s *Server) PostNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) GetNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) PatchNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
