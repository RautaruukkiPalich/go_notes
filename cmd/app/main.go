package main

import (
	"log"

	"github.com/rautaruukkipalich/go_notes/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}