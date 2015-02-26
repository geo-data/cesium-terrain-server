package handlers

import (
	"io"
	"net/http"
)

type multiWriter struct {
	writers []http.ResponseWriter
}

func (t *multiWriter) Header() http.Header {
	for _, w := range t.writers[1:] {
		w.Header()
	}

	return t.writers[0].Header()
}

func (t *multiWriter) WriteHeader(status int) {
	for _, w := range t.writers {
		w.WriteHeader(status)
	}
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

// MultiWriter is inspired by io.MultiWriter
func MultiWriter(writers ...http.ResponseWriter) http.ResponseWriter {
	w := make([]http.ResponseWriter, len(writers))
	copy(w, writers)
	return &multiWriter{w}
}
