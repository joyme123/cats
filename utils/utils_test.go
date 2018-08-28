package utils

import (
	"testing"
	"time"
)

func TestFmtGMT(t *testing.T) {
	loc, err := time.LoadLocation("UTC")

	if err != nil {
		t.Error("loc error")
	}

	time := time.Date(2018, 8, 26, 1, 59, 0, 0, loc)

	gmtTime := FmtGMT(time)

	if gmtTime != "Sun, 26 Aug 2018 01:59:00 GMT" {
		t.Error("time parse error")
	}
}
