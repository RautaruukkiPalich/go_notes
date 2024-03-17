package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rautaruukkipalich/go_notes/internal/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	store  store.Store
	cache store.Cache
	router *mux.Router
	logger *logrus.Logger
}

func NewServer(
	store store.Store,
	cache store.Cache,
) *Server {
	return &Server{
		store: store,
		cache: cache,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	s.router.Handle("/notes", s.GetNotes()).Methods(http.MethodGet)
	s.router.Handle("/notes", s.PostNote()).Methods(http.MethodPost)
	s.router.Handle("/notes/{id}", s.AuthMiddleware(http.Handler(s.GetNote()))).Methods(http.MethodGet)
	s.router.Handle("/notes/{id}", s.PatchNote()).Methods(http.MethodPatch)
	s.router.Handle("/notes/{id}", s.DeleteNote()).Methods(http.MethodDelete)
}