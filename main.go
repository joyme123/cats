package main

import (
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/core/index"
	"github.com/joyme123/cats/core/serveFile"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var conf config.Config

	conf = config.Config{
		Addr:      "127.0.0.1",
		Port:      8089,
		ServeFile: "/home/jiang/go/src/github.com/joyme123/cats/test-web",
		Index:     "index.html"}

	var server http.Server

	server.Config(&conf)

	if conf.Index != "" {
		indexComp := index.Index{File: conf.Index}
		server.Register(&indexComp)
	}

	if conf.ServeFile != "" {
		serveFileComp := serveFile.ServeFile{RootDir: conf.ServeFile}
		server.Register(&serveFileComp)
	}

	server.Start()
}
