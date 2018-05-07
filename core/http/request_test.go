package http

import (
	"testing"
)

func TestParse(t *testing.T) {
	var req Request

	in := []byte("Content-Type: text/html\r\nContent-Length: 9\r\n\r\ntest test")
	req.Parse(in)

	if req.Headers["Content-Type"] != "text/html" {
		t.Errorf("request Parse Headers Error\n%v", req.Headers["Content-Type"])
	}

	if req.Headers["Content-Length"] != "9" {
		t.Errorf("request Parse Headers Error\n%v", req.Headers["Content-Length"])
	}

	if len(req.Body) != 9 {
		t.Errorf("request Parse Body Error\n%v", req.Body)
	}
}
