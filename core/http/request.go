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
}

func (req *Request) logger(out io.Writer) {
	fmt.Fprintf(out, "%v %v %v %v\n", req.Method, req.URI, req.Version, req.Headers)
}

func (req *Request) loggerBody(out io.Writer) {
	fmt.Fprintf(out, "%v", req.Body)
}
