package main

import (
	"fmt"
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/core/index"
	"github.com/joyme123/cats/core/serveFile"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var vhost1 config.VHost

	vhost1 = config.VHost{
		Addr:      "127.0.0.1",
		Port:      8089,
		ServeFile: "/home/jiang/projects/test-web",
		Index:     []string{"index.htm", "index.html"}}

	var conf config.Config

	conf.VHosts = append(conf.VHosts, vhost1)

	for _, vhost := range conf.VHosts {

		fmt.Printf("%v\n", vhost)

		var server http.Server

		server.Context(&vhost)

		if len(vhost.Index) != 0 {
			indexComp := index.Index{}
			indexComp.New(server.GetContext(), &vhost)
			server.Register(&indexComp)
		}

		if vhost.ServeFile != "" {
			serveFileComp := serveFile.ServeFile{}
			serveFileComp.New(server.GetContext(), &vhost)
			server.Register(&serveFileComp)
		}

		server.Start()
	}

}
