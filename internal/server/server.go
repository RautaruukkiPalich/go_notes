package server

import (
	"github.com/gorilla/mux"
	"github.com/rautaruukkipalich/go_notes/internal/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	store  store.Store
	router *mux.Router
	logger *logrus.Logger
}

func NewServer(store store.Store) *Server {
	return &Server{
		store: store,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
}