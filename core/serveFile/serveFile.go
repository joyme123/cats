package serveFile

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/joyme123/cats/utils"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type ServeFile struct {
	http.BaseComponent
	Index   int
	Context *http.VhostContext
}

func (serverFile *ServeFile) New(site *config.Site, context *http.VhostContext) {
	serverFile.Context = context
}

func (serverFile *ServeFile) serveFile(filepath string, req *http.Request, resp *http.Response) {

	// 检查是否有range头
	if rangeValue, ok := req.Headers["range"]; ok {
		// 并且是GET
		if req.Method == "GET" {
			rangeData, err := parseRange([]byte(rangeValue))
			if err == nil && rangeData.unit == "bytes" {
				data, err := ioutil.ReadFile(filepath)
				if err != nil {
					resp.Error404()
				} else {
					dataLen := len(data)
					// 检查最后一个range的范围是否超出
					if rangeData.parts[len(rangeData.parts)-1].end >= dataLen {
						// 超出最大文件长度，返回416
						resp.StatusCode = 416
						resp.Desc = "Requested Range Not Satisfiable"
						return
					} else {
						// 组装multipart的响应进行返回
						lastIndex := strings.LastIndex(filepath, ".")

						var ctype string

						if lastIndex > 0 {
							var ok bool
							if ctype, ok = utils.GetMimeByExt(string([]byte(filepath)[lastIndex+1:])); !ok {

								ctype = "text/plain"
							}
						} else {
							ctype = "text/plain"
						}

						var boundary string
						resp.Body, boundary = mergeMultiRange(data, rangeData.parts, ctype)
						resp.StatusCode = 206
						resp.Desc = "Partial Content"
						resp.AppendHeader("content-type", "multipart/byteranges; boundary="+boundary)
						return
					}
				}
			}

		}
	}

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
	defer serverFile.Next(req, resp)

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
