package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/joyme123/cats/utils"
)

type Response struct {
	Version      string
	StatusCode   int
	Desc         string
	Headers      map[string]string
	HeaderCookie []string
	Body         []byte
	Writer       io.Writer
}

func (resp *Response) Init(writer io.Writer) {
	resp.Writer = writer
	resp.Headers = make(map[string]string)
}

// AppendHeader会将所有的Header处理成小写 向Response中添加响应头
func (resp *Response) AppendHeader(k string, v string) {
	k = strings.ToLower(k)
	if k == "set-cookie" {
		resp.HeaderCookie = append(resp.HeaderCookie, v)
	} else {
		resp.Headers[k] = v
	}

}

// 将响应转为字符串
func (resp *Response) toBytes() []byte {

	// log.Printf("response body的为:%v,%s", len(resp.Body), string(resp.Body))
	var buf bytes.Buffer
	startLine := fmt.Sprintf("%v %v %v\r\n", resp.Version, resp.StatusCode, resp.Desc)
	buf.WriteString(startLine)
	for k, v := range resp.Headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}

	// 写出cookie
	for _, v := range resp.HeaderCookie {
		buf.WriteString("set-cookie: " + v + "\r\n")
	}

	buf.WriteString("content-length: " + strconv.Itoa(len(resp.Body)) + "\r\n")
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

func (resp *Response) Error502() {
	resp.StatusCode = 502
	resp.Desc = "error"
	resp.Body = []byte("bad gateway")
}

func (resp *Response) Status304() {
	resp.StatusCode = 304
	resp.Desc = "Not Modifie"
}

func (resp *Response) out() {

	// 在输出信息的时候，添加Date头

	resp.AppendHeader("date", utils.FmtGMT(time.Now()))

	_, err := resp.Writer.Write(resp.toBytes())

	if err != nil {
		log.Println(err)
	}
}

// 清空状态
func (resp *Response) Clear() {
	resp.Version = ""
	resp.StatusCode = 0
	resp.Desc = ""
	resp.Headers = make(map[string]string)
	resp.Body = nil
}
