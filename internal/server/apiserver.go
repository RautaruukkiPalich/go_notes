package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/store/sqlstore"
)

const (
	bindAddr = "localhost:8088"
	dbUri = "postgres://postgres:postgres@localhost:5432/go_notes?sslmode=disable"
)

func Start() error {
	db, err := newDB(dbUri)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store, err := sqlstore.New(db)
	if err != nil {
		panic(err)
	}

	s := NewServer(store)

	server := &http.Server{
		Addr: bindAddr,
		Handler: s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.logger.Print(
		fmt.Sprintf(
			"server up on '%s'",
			server.Addr,
		),
	)

	return server.ListenAndServe()
}

func newDB(databaseURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURI)
	
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}