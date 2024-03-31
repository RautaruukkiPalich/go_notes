package sqlstore_test

import (
	"os"
	"testing"
)

const (
	databaseUrl string = "postgres://postgres:postgres@localhost:5432/go_notes_test?sslmode=disable"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}