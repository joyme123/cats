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

func (serverFile *ServeFile) New(site *config.Site, context *http.Context) {
	serverFile.RootDir = site.Root
	serverFile.Context = context
}

func (serverFile *ServeFile) serveFile(filepath string) {

	var fileerr error
	serverFile.resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		serverFile.resp.Error404()
	} else {
		serverFile.resp.StatusCode = 200
		serverFile.resp.Desc = "OK"
	}
}

func (serverFile *ServeFile) Start() {

}

func (serverFile *ServeFile) commonHeaders() {
	serverFile.resp.AppendHeader("connection", "keep-alive")
	serverFile.resp.AppendHeader("server", "cats")
}

func (serverFile *ServeFile) Serve(req *http.Request, resp *http.Response) {
	serverFile.req = req
	serverFile.resp = resp

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

	filepath = serverFile.RootDir + filepath

	// 文件夹结尾,自动加上index文件
	if strings.HasSuffix(filepath, "/") {

		if indexFiles, ok := serverFile.Context.KeyValue["IndexFiles"]; ok {
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

	serverFile.req.Context["FilePath"] = filepath

	log.Println("server file:", filepath)
	serverFile.commonHeaders()
	serverFile.serveFile(filepath)
}

func (serverFile *ServeFile) Shutdown() {

}

func (serverFile *ServeFile) GetIndex() int {

	return serverFile.Index
}
