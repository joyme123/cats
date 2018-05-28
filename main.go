package main

import (
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var conf config.Config

	conf = config.Config{
		Addr:    "127.0.0.1",
		Port:    8089,
		RootDir: "/home/jiang/opensource/go/bin/src/github.com/joyme123/cats/test-web",
		Index:   "index.html"}

	var server http.Server

	server.Config(&conf)
	server.Start()
}
