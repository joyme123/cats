package cache

import (
	"strings"
	"testing"
	"time"
)

func TestEtag(t *testing.T) {
	etag := Etag("./doc.md")

	if strings.Index(etag, "-") != 8 {
		t.Error("etag error")
	}
}

func TestFmtGMT(t *testing.T) {
	loc, err := time.LoadLocation("UTC")

	if err != nil {
		t.Error("loc error")
	}

	time := time.Date(2018, 8, 26, 1, 59, 0, 0, loc)

	gmtTime := fmtGMT(time)

	if gmtTime != "Sun, 26 Aug 2018 01:59:00 GMT" {
		t.Error("time parse error")
	}
}

func TestLastModified(t *testing.T) {
	lm := LastModified("doc.md")

	if len(lm) != 29 {
		t.Error("last modified time error")
	}
}
