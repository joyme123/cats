package cache

import (
	"strings"
	"testing"
	"time"
)

func TestEtag(t *testing.T) {
	etag := Etag("./cache.go")

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
	lm := LastModified("./cache.go")

	if len(lm) != 29 {
		t.Error("last modified time error")
	}
}

func TestParseGMT(t *testing.T) {
	gmtTime, _ := ParseGMT("Sun, 26 Aug 2018 01:59:00 GMT")

	if gmtTime != time.Date(2018, time.August, 26, 1, 59, 0, 0, time.UTC) {
		t.Error("parse GMT time error")
	}

	gmtTime, err := ParseGMT("26 Aug 2018 01:59:00 GMT")

	if err == nil {
		t.Error("parse GMT time error")
	}

}

func TestCompareFileModifiedTime(t *testing.T) {
	if !CompareFileModifiedTime("./cache.go", "Sun, 26 Aug 2018 01:59:00 GMT") {
		t.Error("compare file modifired time error")
	}

	if CompareFileModifiedTime("./cache.go", "Sun, 26 Aug 2090 01:59:00 GMT") {
		t.Error("compare file modifired time error")
	}
}
