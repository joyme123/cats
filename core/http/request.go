package http

import (
	"bytes"
	"log"
)

type Request struct {
	Headers map[string]string
	Body    []byte
}

// Parse 函数用来解析输入流来构造请求头和请求体
func (req *Request) Parse(in []byte) {

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	var buf bytes.Buffer
	var k string
	var v string
	var last byte

	endline := false //是否匹配完一行
	kend := false    //当前行是否已经匹配到: 了

	// 一个简单的有限状态机
	for index, c := range in {
		switch c {
		case '\r':
			if last != '\n' { //不是连续两次\r\n
				v = buf.String()
				buf.Reset()
				req.Headers[k] = v
			}

			last = c
		case '\n':
			if last == '\r' { //上一次是\r,则代表一行的结束
				if endline { // 处于endline状态,header解析完毕
					req.Body = in[index+1:] //构造body
					break                   // end for
				}
				endline = true
				kend = false
			} else {
				log.Fatal("unexpect \\n")
			}
			last = c
		case ':':
			if !kend {
				k = buf.String()
				buf.Reset()
				kend = true
			}
			last = c
		case ' ':
			last = c
			endline = false
		default:
			buf.WriteByte(c)
			last = c
			endline = false
		}
	}
}
