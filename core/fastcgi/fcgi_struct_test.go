package fastcgi

import "testing"

func TestFCGIHeaderToBlob(t *testing.T) {
	var header FCGIHeader

	header.Version = 1
	header.Type = 1
	header.RequestID = 1
	header.ContentLength = 8
	header.PaddingLength = 0
	header.Reserved = 0

	headerBytes := header.ToBlob()

	expectVal := []byte{1, 1, 0, 1, 0, 8, 0, 0}

	// equals to for i,_ := range ...
	for i := range expectVal {
		if headerBytes[i] != expectVal[i] {
			t.Errorf("header byte index %d error, actual: %v, expect: %v", i, headerBytes[i], expectVal[i])
		}
	}
}

func TestFCGIBeginRequestRecordToBlob(t *testing.T) {
	var record FCGIBeginRequestRecord

	var header FCGIHeader

	header.Version = 1
	header.Type = 1
	header.RequestID = 1
	header.ContentLength = 8
	header.PaddingLength = 0
	header.Reserved = 0

	record.Header = header

	var body FCGIBeginRequestBody

	body.Role = 1
	body.Flags = 0

	record.Body = body

	blob := record.ToBlob()

	expect := []byte{1, 1, 0, 1, 0, 8, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}

	for i := range expect {
		if blob[i] != expect[i] {
			t.Errorf("record byte index %d error, actual: %v, expect: %v", i, blob[i], expect[i])
		}
	}
}

func TestFCGIEndRequestRecordToBlob(t *testing.T) {
	var record FCGIEndRequestRecord

	var header FCGIHeader

	header.Version = 1
	header.Type = 1
	header.RequestID = 1
	header.ContentLength = 8
	header.PaddingLength = 0
	header.Reserved = 0

	record.Header = header

	var body FCGIEndRequestBody

	body.AppStatus = 255
	body.ProtocolStatus = 0

	record.Body = body

	blob := record.ToBlob()

	expect := []byte{1, 1, 0, 1, 0, 8, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0}

	for i := range expect {
		if blob[i] != expect[i] {
			t.Errorf("record byte index %d error, actual: %v, expect: %v", i, blob[i], expect[i])
		}
	}
}

func TestFCGIParamsRecordToBlob(t *testing.T) {
	var record FCGIParamsRecord

	var header FCGIHeader

	header.Version = 1
	header.Type = 1
	header.RequestID = 1
	header.ContentLength = 23
	header.PaddingLength = 0
	header.Reserved = 0

	record.Header = header

	var body FCGINameValuePair11

	body.NameLength = 12
	body.NameData = []byte("Content-Type")
	body.ValueLength = 9
	body.ValueData = []byte("text/html")

	// 这里之所以要用取地址符，是因为
	// https://stackoverflow.com/questions/40823315/go-x-does-not-implement-y-method-has-a-pointer-receiver
	record.Body = &body

	blob := record.ToBlob()

	expect := []byte{1, 1, 0, 1, 0, 23, 0, 0, 12, 9, 'C', 'o', 'n', 't', 'e', 'n', 't', '-', 'T', 'y', 'p', 'e', 't', 'e', 'x', 't', '/', 'h', 't', 'm', 'l'}

	for i := range expect {
		if blob[i] != expect[i] {
			t.Errorf("record byte index %d error, actual: %v, expect: %v", i, blob[i], expect[i])
		}
	}
}

func TestFCGIStdinRecordToBlob(t *testing.T) {
	var record FCGIStdinRecord

	var header FCGIHeader

	header.Version = 1
	header.Type = 1
	header.RequestID = 1
	header.ContentLength = 23
	header.PaddingLength = 0
	header.Reserved = 0

	record.Header = header
	record.Body = []byte("abcdabcdabcdabcd")

	blob := record.ToBlob()

	expect := []byte{1, 1, 0, 1, 0, 23, 0, 0, 'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd', 'a', 'b', 'c', 'd'}

	for i := range expect {
		if blob[i] != expect[i] {
			t.Errorf("record byte index %d error, actual: %v, expect: %v", i, blob[i], expect[i])
		}
	}
}
