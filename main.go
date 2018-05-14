package main

import (
	"log"

	"github.com/joyme123/cats/core/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var server http.Server

	server.Start()
}
