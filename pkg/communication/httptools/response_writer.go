package httptools

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

// A ResponseWriter is used to capture set status codes.
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker // need this for sentry
	http.Flusher  // need this for sentry
	io.ReaderFrom // need this for sentry
	Status() int
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

// Flush sends any buffered data to the client.
func (w *responseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		panic(fmt.Errorf("ResponseWriter doesn't implement http.Flusher"))
	}

	flusher.Flush()
}

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
func (w *responseWriter) ReadFrom(r io.Reader) (int64, error) {
	reader, ok := w.ResponseWriter.(io.ReaderFrom)
	if !ok {
		panic(fmt.Errorf("ResponseWriter doesn't implement io.ReaderFrom"))
	}

	return reader.ReadFrom(r)
}

// NewResponseWriter returns a new [ResponseWriter].
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{w, -1}
}

// WriteHeader sets the internal status value of a [ResponseWriter].
func (w *responseWriter) WriteHeader(status int) {
	if w.status != -1 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Status returns the status code of a [ResponseWriter].
func (w responseWriter) Status() int {
	return w.status
}

// Hijack lets the caller take over the connection.
// After a call to Hijack the HTTP server library
// will not do anything else with the connection.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}
