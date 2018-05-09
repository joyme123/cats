package http

import (
	"testing"
)

func TestParseHeader(t *testing.T) {
	var req Request

	in := []byte("GET /index.html HTTP/1.1\r\nContent-Type: text/html\r\nContent-Length: 9\r\n\r\n")
	req.ParseHeader(in)

	if req.Method != "GET" {
		t.Errorf("request Parse Method Error\n%v", req.Method)
	}

	if req.URL != "/index.html" {
		t.Errorf("request Parse URL Error\n%v", req.URL)
	}

	if req.Version != "HTTP/1.1" {
		t.Errorf("request Parse Version Error\n%v", req.Version)
	}

	if req.Headers["Content-Type"] != "text/html" {
		t.Errorf("request Parse Headers Error\n%v", req.Headers["Content-Type"])
	}

	if req.Headers["Content-Length"] != "9" {
		t.Errorf("request Parse Headers Error\n%v", req.Headers["Content-Length"])
	}
}
