package serveFile

import (
	"io/ioutil"
	"log"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type ServeFile struct {
	Index   int
	Context *http.VhostContext
}

func (serverFile *ServeFile) New(site *config.Site, context *http.VhostContext) {
	serverFile.Context = context
}

func (serverFile *ServeFile) serveFile(filepath string, req *http.Request, resp *http.Response) {

	var fileerr error
	resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		resp.Error404()
	} else {
		resp.StatusCode = 200
		resp.Desc = "OK"
	}
}

func (serverFile *ServeFile) Start() {

}

func (serverFile *ServeFile) commonHeaders(resp *http.Response) {
	resp.AppendHeader("connection", "keep-alive")
	resp.AppendHeader("server", "cats")
}

func (serverFile *ServeFile) Serve(req *http.Request, resp *http.Response) {

	if resp.StatusCode != 0 {
		return
	}

	filepath, ok := req.Context["FilePath"]

	if !ok {
		log.Println("serve file error: not found filepath in request context")
		resp.Error404()
		return
	}

	if str, ok := filepath.(string); ok {
		log.Println("server file:", filepath)
		serverFile.commonHeaders(resp)
		serverFile.serveFile(str, req, resp)
	} else {
		resp.Error404()
		return
	}

}

func (serverFile *ServeFile) Shutdown() {

}

func (serverFile *ServeFile) GetIndex() int {

	return serverFile.Index
}

func (serverFile *ServeFile) GetContainer() string {
	return "location"
}
