package index

import (
	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Index struct {
	Files      []string
	Index      int
	IndexFiles []string
	Context    *http.Context // handler的context
	req        *http.Request
	resp       *http.Response
}

func (index *Index) New(vhost *config.VHost) {
	index.IndexFiles = vhost.Index
}

func (index *Index) Start(context *http.Context) {
	index.Context = context
	index.Context.KeyValue["IndexFiles"] = index.IndexFiles
}

func (index *Index) Serve(req *http.Request, resp *http.Response) {

}

func (index *Index) Shutdown() {

}

func (index *Index) GetIndex() int {

	return index.Index
}
