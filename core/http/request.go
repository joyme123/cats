package http

import (
	"fmt"
	"io"
)

// header field一律小写存储
type Request struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    []byte
	Context map[string]interface{}
}

func (req *Request) logger(out io.Writer) {
	fmt.Fprintf(out, "%v %v %v %v\n", req.Method, req.URI, req.Version, req.Headers)
}

func (req *Request) loggerBody(out io.Writer) {
	fmt.Fprintf(out, "%v", req.Body)
}

func (req *Request) Clear() {
	req.Headers = make(map[string]string)
}
