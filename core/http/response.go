package http

import "bytes"

type Response struct {
	Headers map[string]string
	Body    []byte
}

// 向Response中添加响应头
func (resp *Response) appendHeader(k string, v string) {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	resp.Headers[k] = v
}

// 将响应转为字符串
func (resp *Response) toString() string {
	var buf bytes.Buffer
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	for k, v := range resp.Headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}
	buf.WriteString("\r\n")
	buf.Write(resp.Body)

	return buf.String()
}
