package main

import "github.com/rautaruukkipalich/go_notes/internal/server"

func main() {
	err := server.Start()

	panic(err)
}