package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rautaruukkipalich/go_notes/internal/cachestore/rediscache"
	"github.com/rautaruukkipalich/go_notes/internal/store/sqlstore"
)

const (
	bindAddr = "localhost:8088"
	dbUri = "postgres://postgres:postgres@localhost:5432/go_notes?sslmode=disable"
	cacheUri = "redis://user:password@localhost:6379/0?protocol=3"
)

func Start() error {
	db, err := newDB(dbUri)
	if err != nil {
		return err
	}
	defer db.Close()

	store, err := sqlstore.New(db)
	if err != nil {
		return err
	}

	cacheDB, err := newCache(cacheUri)
	if err != nil {
		return err
	}
	defer cacheDB.Close()

	cache, err := rediscache.New(cacheDB)
	if err != nil {
		return err
	}

	s := NewServer(store, cache)
	s.configureRouter()

	server := &http.Server{
		Addr: bindAddr,
		Handler: s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.logger.Info(
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

func newCache(redisUri string) (*redis.Client, error) {
    opts, err := redis.ParseURL(redisUri)
    if err != nil {
        return nil, err
    }

    return redis.NewClient(opts), nil
}