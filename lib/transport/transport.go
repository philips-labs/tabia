package transport

import (
	"fmt"
	"io"
	"net/http"
)

// TeeRoundTripper copies request bodies to stdout.
type TeeRoundTripper struct {
	http.RoundTripper
	Writer io.Writer
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
func (t TeeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Fprintf(t.Writer, "%s: %s ", req.Method, req.URL)
	if req.Body != nil {
		req.Body = struct {
			io.Reader
			io.Closer
		}{
			Reader: io.TeeReader(req.Body, t.Writer),
			Closer: req.Body,
		}
	}

	return t.RoundTripper.RoundTrip(req)
}
