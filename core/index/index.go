package hub

import (
	"net/url"
	"strings"

	"github.com/joyme123/cats/core/http"
)

type Index struct {
	Files []string
	Index int
	req   *http.Request
	resp  *http.Response
}

func (index *Index) Start() {

}

func (index *Index) Serve(req *http.Request, resp *http.Response) {
	index.req = req
	index.resp = resp

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

	// 文件夹结尾,自动加上index文件
	if strings.HasSuffix(filepath, "/") {
		filepath = filepath + "index.html"
	}

	req.URI = filepath
}

func (index *Index) Shutdown() {

}

func (index *Index) GetIndex() int {

	return index.Index
}
