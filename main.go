package main

import (
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/core/index"
	"github.com/joyme123/cats/core/mime"
	"github.com/joyme123/cats/core/serveFile"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var conf config.Config

	conf = config.Config{
		Addr:      "127.0.0.1",
		Port:      8089,
		ServeFile: "/home/jiang/projects/test-web",
		Index:     []string{"index.htm", "index.html"}}

	var server http.Server

	server.Config(&conf)

	if len(conf.Index) != 0 {
		indexComp := index.Index{}
		indexComp.New(&conf)
		server.Register(&indexComp)
	}

	if conf.ServeFile != "" {
		serveFileComp := serveFile.ServeFile{}
		serveFileComp.New(&conf)
		server.Register(&serveFileComp)
	}

	mimeComp := mime.Mime{}
	mimeComp.New(&conf)
	server.Register(&mimeComp)

	server.Start()
}
