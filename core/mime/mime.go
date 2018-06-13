package mime

import (
	"strings"

	"github.com/joyme123/cats/config"
	"github.com/joyme123/cats/core/http"
)

type Mime struct {
	mimes   map[string]string
	Context *http.Context
	req     *http.Request
	resp    *http.Response
	Index   int
}

func (mime *Mime) New(config *config.Config) {
	mime.mimes = make(map[string]string)

	// RFC https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Complete_list_of_MIME_types
	mime.mimes["aac"] = "audio/aac"
	mime.mimes["abw"] = "application/x-abiword"
	mime.mimes["arc"] = "application/octet-stream"
	mime.mimes["avi"] = "video/x-msvideo"
	mime.mimes["azw"] = "application/vnd.amazon.ebook"
	mime.mimes["bin"] = "application/octet-stream"
	mime.mimes["bz"] = "application/x-bzip"
	mime.mimes["bz2"] = "application/x-bzip2"
	mime.mimes["csh"] = "application/x-csh"
	mime.mimes["css"] = "text/css"
	mime.mimes["csv"] = "text/csv"
	mime.mimes["doc"] = "application/msword"
	mime.mimes["epub"] = "application/epub+zip"
	mime.mimes["gif"] = "image/gif"
	mime.mimes["htm"] = "text/html"
	mime.mimes["html"] = "text/html"
	mime.mimes["ico"] = "image/x-icon"
	mime.mimes["ics"] = "text/calendar"
	mime.mimes["jar"] = "application/java-archive"
	mime.mimes["jpeg"] = "image/jpeg"
	mime.mimes["jpg"] = "image/jpeg"
	mime.mimes["js"] = "application/javascript"
	mime.mimes["json"] = "application/json"
	mime.mimes["mid"] = "audio/midi"
	mime.mimes["midi"] = "audio/midi"
	mime.mimes["mpeg"] = "video/mpeg"
	mime.mimes["mpkg"] = "application/vnd.apple.installer+xml"
	mime.mimes["odp"] = "application/vnd.oasis.opendocument.presentation"
	mime.mimes["ods"] = "application/vnd.oasis.opendocument.spreadsheet"
	mime.mimes["odt"] = "application/vnd.oasis.opendocument.text"
	mime.mimes["oga"] = "audio/ogg"
	mime.mimes["ogv"] = "video/ogg"
	mime.mimes["ogx"] = "application/ogg"
	mime.mimes["pdf"] = "application/pdf"
	mime.mimes["ppt"] = "application/vnd.ms-powerpoint"
	mime.mimes["rar"] = "application/x-rar-compressed"
	mime.mimes["rtf"] = "application/rtf"
	mime.mimes["sh"] = "application/x-sh"
	mime.mimes["svg"] = "image/svg+xml"
	mime.mimes["swf"] = "application/x-shockwave-flash"
	mime.mimes["tar"] = "application/x-tar"
	mime.mimes["tif"] = "image/tiff"
	mime.mimes["tiff"] = "image/tiff"
	mime.mimes["ttf"] = "application/x-font-ttf"
	mime.mimes["vsd"] = "application/vnd.visio"
	mime.mimes["wav"] = "audio/x-wav"
	mime.mimes["weba"] = "audio/webm"
	mime.mimes["webm"] = "video/webm"
	mime.mimes["webp"] = "image/webp"
	mime.mimes["woff"] = "application/x-font-woff"
	mime.mimes["xhtml"] = "application/xhtml+xml"
	mime.mimes["xls"] = "application/vnd.ms-excel"
	mime.mimes["xml"] = "application/xml"
	mime.mimes["xul"] = "application/vnd.mozilla.xul+xml"
	mime.mimes["zip"] = "application/zip"
	mime.mimes["3gp"] = "video/3gpp"
	mime.mimes["3g2"] = "video/3gpp2"
	mime.mimes["7z"] = "application/x-7z-compressed"
}

func (mime *Mime) Start(context *http.Context) {
	mime.Context = context
}

func (mime *Mime) Serve(req *http.Request, resp *http.Response) {
	mime.req = req
	mime.resp = resp
	ctype, haveType := mime.resp.Headers["content-type"]

	if !haveType {
		filepath := mime.Context.KeyValue["FilePath"].(string)
		lastIndex := strings.LastIndex(filepath, ".")

		if lastIndex > 0 {

			var ok bool
			if ctype, ok = mime.mimes[string([]byte(filepath)[lastIndex+1:])]; !ok {

				ctype = "text/plain"
			}
		} else {
			ctype = "text/plain"
		}
		mime.resp.Headers["content-type"] = ctype
	}
}

func (mime *Mime) Shutdown() {

}

func (mime *Mime) GetIndex() int {

	return mime.Index
}
