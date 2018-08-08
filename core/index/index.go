package index

import (
	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Index struct {
	Files      []string
	Index      int
	IndexFiles []string
	Context    *http.VhostContext // vhost的context
}

func (index *Index) New(site *config.Site, context *http.VhostContext) {
	index.IndexFiles = site.Index
	index.Context = context
	index.Context.KeyValue["IndexFiles"] = index.IndexFiles
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

func (index *Index) GetContainer() string {
	return "vhost"
}
