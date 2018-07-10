package index

import (
	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Index struct {
	Files   []string
	Index   int
	Context *http.Context
	req     *http.Request
	resp    *http.Response
}

func (index *Index) New(context *http.Context, vhost *config.VHost) {
	index.Context = context
	index.Context.KeyValue["IndexFiles"] = vhost.Index
}

func (index *Index) Start() {

}

func (index *Index) Serve(req *http.Request, resp *http.Response) {

}

func (index *Index) Shutdown() {

}

func (index *Index) GetIndex() int {

	return index.Index
}
