package main

import (
	"log"

	"github.com/chrollo-lucider-12/proxy/server"
)

func main() {

	if err := server.Run(); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
