package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"strconv"
	"strings"
)

type Response struct {
	Version    string
	StatusCode int
	Desc       string
	Headers    map[string]string
	Body       []byte
	Writer     io.Writer
}

// 向Response中添加响应头
func (resp *Response) appendHeader(k string, v string) {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	resp.Headers[k] = v
}

// 将响应转为字符串
func (resp *Response) toBytes() []byte {
	var buf bytes.Buffer
	startLine := fmt.Sprintf("%v %v %v\r\n", resp.Version, resp.StatusCode, resp.Desc)
	buf.WriteString(startLine)
	for k, v := range resp.Headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}
	buf.WriteString("Content-Length: " + strconv.Itoa(len(resp.Body)) + "\r\n")
	buf.WriteString("\r\n")
	buf.Write(resp.Body)

	return buf.Bytes()
}

func (resp *Response) commonHeaders() {
	resp.appendHeader("Connection", "keep-alive")
	resp.appendHeader("server", "cats")
}

func (resp *Response) serveFile(filepath string) {
	resp.commonHeaders()

	var fileerr error
	resp.Body, fileerr = ioutil.ReadFile(filepath)
	if fileerr != nil {
		resp.Error404()
	} else {
		resp.StatusCode = 200
		resp.Desc = "OK"

		ctype, haveType := resp.Headers["Content-Type"]
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
			resp.Headers["Content-Type"] = ctype
		}
	}

	resp.out()
}

// 只有serveFile调用
func (resp *Response) Error404() {
	//TODO: 读取不到文件,暂时返回404
	resp.StatusCode = 404
	resp.Desc = "error"
	resp.Body = []byte("page not found")
}

func (resp *Response) Error400() {
	resp.commonHeaders()
	resp.StatusCode = 400
	resp.Desc = "error"
	resp.Body = []byte("bad request")
}

func (resp *Response) out() {
	_, err := resp.Writer.Write(resp.toBytes())

	if err != nil {
		log.Println(err)
	}
}
