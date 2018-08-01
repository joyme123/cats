package serveFile

import (
	"io/ioutil"
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
	"github.com/joyme123/cats/utils"
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

	if indexFiles, ok := serverFile.Context.KeyValue["IndexFiles"]; ok {
		filepath = utils.GetAbsolutePath(serverFile.RootDir, req.URI, indexFiles.([]string))
	} else {
		filepath = utils.GetAbsolutePath(serverFile.RootDir, req.URI, make([]string, 0, 0))
	}

	if filepath == "400" {
		resp.Error400()
		return
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

func (serverFile *ServeFile) GetContainer() string {
	return "location"
}
