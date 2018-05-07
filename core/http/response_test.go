package http

import "testing"

func TestAppendHeader(t *testing.T) {
	var resp Response

	resp.appendHeader("Content-Type", "text/html")
	if resp.Headers["Content-Type"] != "text/html" {
		t.Error("appendHeader error")
	}
}

func TestToString(t *testing.T) {
	var resp Response

	resp.appendHeader("Content-Type", "text/html")
	resp.Body = []byte("test test")

	str := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\ntest test"

	if resp.toString() != str {
		t.Errorf("response construct error\n%v", resp.toString())
	}
}
