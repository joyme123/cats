package serveFile

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type ServeFile struct {
	RootDir string
	Index   int
	Context *http.Context
	req     *http.Request
	resp    *http.Response
}

func (server *ServeFile) New(config *config.Config) {
	server.RootDir = config.ServeFile
}

func (server *ServeFile) serveFile(filepath string) {

	var fileerr error
	server.resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		server.resp.Error404()
	} else {
		server.resp.StatusCode = 200
		server.resp.Desc = "OK"
	}
}

func (server *ServeFile) Start(context *http.Context) {
	server.Context = context
}

func (server *ServeFile) commonHeaders() {
	server.resp.AppendHeader("Connection", "keep-alive")
	server.resp.AppendHeader("server", "cats")
}

func (server *ServeFile) Serve(req *http.Request, resp *http.Response) {
	server.req = req
	server.resp = resp

	var filepath string
	if strings.HasPrefix(req.URI, "http") {
		u, err := url.Parse(req.URI)
		if err != nil {
			resp.Error400()
			return
		} else {
			filepath = u.Path
		}
	} else {
		filepath = req.URI
	}

	filepath = server.RootDir + filepath

	// 文件夹结尾,自动加上index文件
	if strings.HasSuffix(filepath, "/") {

		if indexFiles, ok := server.Context.KeyValue["IndexFiles"]; ok {
			for _, v := range indexFiles.([]string) {
				_, err := os.Stat(filepath + v)
				if err == nil {
					// 文件存在
					filepath = filepath + v
					break
				}
			}
		} else {
			// 默认为index.html
			filepath = filepath + "index.html"
		}

	}

	server.Context.KeyValue["FilePath"] = filepath

	log.Println("server file:", filepath)
	server.commonHeaders()
	server.serveFile(filepath)
}

func (server *ServeFile) Shutdown() {

}

func (server *ServeFile) GetIndex() int {

	return server.Index
}
