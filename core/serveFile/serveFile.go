package hub

import (
	"io/ioutil"
	"log"
	"mime"
	"strings"

	"github.com/joyme123/cats/core/http"
)

type ServeFile struct {
	RootDir string
	Index   int
	req     *http.Request
	resp    *http.Response
}

func (root *ServeFile) serveFile(filepath string) {

	var fileerr error
	root.resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		root.resp.Error404()
	} else {
		root.resp.StatusCode = 200
		root.resp.Desc = "OK"

		ctype, haveType := root.resp.Headers["Content-Type"]
		if !haveType {
			lastIndex := strings.LastIndex(filepath, ".")
			if lastIndex > 0 {

				//TODO: 根据扩展名做的mime types是不对的,比如js文件会解析为image/jpeg, css会解析为text/html
				ctype = mime.TypeByExtension(string([]byte(filepath)[lastIndex:]))
				log.Println(string([]byte(filepath)[lastIndex:]))
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
			root.resp.Headers["Content-Type"] = ctype
		}
	}
}

func (root *ServeFile) Start() {

}

func (server *ServeFile) Serve(req *http.Request, resp *http.Response) {
	server.req = req
	server.resp = resp

	filepath := server.RootDir + req.URI

	log.Println(filepath)
	server.serveFile(filepath)
}

func (server *ServeFile) Shutdown() {

}

func (server *ServeFile) GetIndex() int {

	return server.Index
}
