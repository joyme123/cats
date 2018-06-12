package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
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
func (resp *Response) AppendHeader(k string, v string) {
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

// 只有serveFile调用
func (resp *Response) Error404() {
	//TODO: 读取不到文件,暂时返回404
	resp.StatusCode = 404
	resp.Desc = "error"
	resp.Body = []byte("page not found")
}

func (resp *Response) Error400() {
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
