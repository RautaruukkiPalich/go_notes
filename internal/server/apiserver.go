package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rautaruukkipalich/go_notes/internal/store/cachestore"
	"github.com/rautaruukkipalich/go_notes/internal/store/sqlstore"
	"github.com/redis/go-redis/v9"
)

const (
	// local app
	// bindAddr = "localhost:8088"
	// dbUri = "postgres://postgres:postgres@localhost:5434/go_notes?sslmode=disable"
	// cacheUri = "redis://localhost:6381"


	// docker app
	bindAddr = "0.0.0.0:8088"
	dbUri = "postgres://postgres:postgres@database:5432/go_notes?sslmode=disable"
	cacheUri = "redis://redis:6379"
)

func Start() error {

	//connect db
	db, err := newDB(dbUri)
	if err != nil {
		return err
	}
	defer db.Close()

	store, err := sqlstore.New(db)
	if err != nil {
		return err
	}

	//connect redis
	redisDB, err := newRedis(cacheUri)
	if err != nil {
		return err
	}
	defer redisDB.Close()

	cache, err := cachestore.New(redisDB)
	if err != nil {
		return err
	}

	s := NewServer(store, cache)

	if err := migrateTables(db); err != nil {
		s.logger.Error(err)
	}

	s.configureRouter()

	if err := s.heatCache(); err != nil {
		return err
	}

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

func newRedis(redisUri string) (*redis.Client, error) {
    opts, err := redis.ParseURL(redisUri)
    if err != nil {
        return nil, err
    }
	cache := redis.NewClient(opts)
	ctx := context.Background()
	if err := cache.Ping(ctx).Err(); err != nil {
		return nil, err
	}

    return cache, nil
}

func migrateTables(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", 
		driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}