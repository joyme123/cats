package serveFile

import (
	"io/ioutil"
	"log"
	"mime"
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

func (server *ServeFile) New(context *http.Context, vhost *config.VHost) {
	server.Context = context
	server.RootDir = vhost.ServeFile
}

func (server *ServeFile) serveFile(filepath string) {

	var fileerr error
	server.resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		server.resp.Error404()
	} else {
		server.resp.StatusCode = 200
		server.resp.Desc = "OK"

		ctype, haveType := server.resp.Headers["Content-Type"]
		if !haveType {
			lastIndex := strings.LastIndex(filepath, ".")
			if lastIndex > 0 {

				//TODO: 根据扩展名做的mime types是不对的,比如js文件会解析为image/jpeg, css会解析为text/html
				ctype = mime.TypeByExtension(string([]byte(filepath)[lastIndex:]))
			}
			if ctype == "" {
				//TODO: 根据文件内容判断文件类型
				// var buf [30]byte
				// file, _ := os.Open(filepath)
				// n, _ := io.ReadFull(file, buf[:])
				// ctype = DetectContentType(buf[:n])
				// _, err := content.Seek(0, io.SeekStart) // rewind to output whole file
				// if err != nil {
				// 	Error(w, "seeker can't seek", StatusInternalServerError)
				// 	return
				// }
				ctype = "*/*"
			}
			log.Printf("%v\n", server.resp.Headers)
			server.resp.Headers["Content-Type"] = ctype
		}
	}
}

func (server *ServeFile) Start() {

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

	log.Println("server file:", filepath)
	server.commonHeaders()
	server.serveFile(filepath)
}

func (server *ServeFile) Shutdown() {

}

func (server *ServeFile) GetIndex() int {

	return server.Index
}
