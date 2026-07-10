//go:build linux

package main

import (
	"log"

	"github.com/Aditramesh/Redis-implementation/server"
)

func main() {
	log.Println("starting the server")
	server.RunTCPASyncServer()
}
