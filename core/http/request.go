package http

type Request struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    []byte
}
