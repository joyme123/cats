package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/fastcgi"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/core/index"
	"github.com/joyme123/cats/core/location"
	"github.com/joyme123/cats/core/mime"
	"github.com/joyme123/cats/core/serveFile"
)

func startServe(sites []config.Site) {
	var server http.Server

	for _, site := range sites {
		server.AddSite(&site)
	}

	server.Config(sites[0].Addr, sites[0].Port)
	server.Init()

	// 根据sites实例化VirtualHost
	for _, site := range sites {
		var vh http.VirtualHost
		vh.Init()
		if len(site.Index) != 0 {
			indexComp := index.Index{}
			indexComp.New(&site, vh.GetContext())
			vh.Register(&indexComp)
		}

		// TODO: 这里先写死location的注册机制
		locationComp := location.Location{}
		locationComp.New(&site, vh.GetContext())
		vh.Register(locationComp)

		// serveFile应该注册到location组件中
		if site.Root != "" {
			serveFileComp := serveFile.ServeFile{}
			serveFileComp.New(&site, vh.GetContext())
			locationComp.Register("/", "", &serveFileComp)
		}

		// fcgipass 应该注册到location组件中
		if site.FCGIPass != "" {
			// 初始化fcgi
			fastcgiComp := fastcgi.FastCGI{}
			fastcgiComp.New(&site, vh.GetContext())
			locationComp.Register("~*", "^(.+)\\.php$", &fastcgiComp)
		}

		mimeComp := mime.Mime{}
		mimeComp.New(&site, vh.GetContext())
		vh.Register(&mimeComp)

		server.SetVirtualHost(site.ServerName, &vh)
	}

	server.Start()
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var site config.Site

	site = config.Site{
		Addr:       "127.0.0.1",
		Port:       8089,
		ServerName: "mysite.com",
		Root:       "/home/jiang/projects/test-web/php",
		Index:      []string{"index.php", "index.html"},
		FCGIPass:   "127.0.0.1:9000"}

	var site2 config.Site

	site2 = config.Site{
		Addr:       "127.0.0.1",
		Port:       8090,
		ServerName: "mysite.com",
		Root:       "/home/jiang/projects/test-web",
		Index:      []string{"index.htm", "index.html"}}

	var site3 config.Site

	site3 = config.Site{
		Addr:       "127.0.0.1",
		Port:       8090,
		ServerName: "mysite2.com",
		Root:       "/home/jiang/projects/test-web/about",
		Index:      []string{"index.htm", "index.html"}}

	var conf config.Config

	conf.Sites = append(conf.Sites, site)
	conf.Sites = append(conf.Sites, site2)
	conf.Sites = append(conf.Sites, site3)

	var wg sync.WaitGroup
	wg.Add(2)
	var groups map[string][]config.Site
	groups = make(map[string][]config.Site)
	for _, site := range conf.Sites {
		// 这里要对site进行整理，拥有同样ip:host的形成一个group
		addrPort := site.Addr + strconv.Itoa(site.Port)
		if _, ok := groups[addrPort]; ok {
			groups[addrPort] = append(groups[addrPort], site)
		} else {
			arr := []config.Site{site}
			groups[addrPort] = arr
		}
	}

	for _, v := range groups {
		go startServe(v)
	}

	// 根据不同的group示例化不同的server,并注入host
	wg.Wait()
}
