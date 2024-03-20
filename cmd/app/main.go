package main

import (
	"log"

	"github.com/rautaruukkipalich/go_notes/internal/server"
)


//	@title			Swagger Example API
//	@version		0.0.1
//	@description	This is a sample Note service

//	@host		localhost:8088
//	@BasePath	/
func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}