package handlers

import (
	"net/http"
)

// ResponseLimiter is an implementation of ResponseLimiter that wraps a
// http.ResponseWriter and limits how much is written to it. If an attempt is
// made to write more bytes than the limit it silently fails.
type LimitResponse struct {
	Writer   http.ResponseWriter
	Limit    Bytes // the maximum number of bytes that can be written to Writer
	written  Bytes // the number of bytes already sent to Writer
	exceeded bool  // has the limit been exceeded?
}

// NewRecorder returns an initialized LimitResponse. This implements the
// LimiterFactory function type.
func NewLimit(writer http.ResponseWriter, limit Bytes) ResponseLimiter {
	return &LimitResponse{
		Writer: writer,
		Limit:  limit,
	}
}

func (this *LimitResponse) LimitExceeded() bool {
	return this.exceeded
}

// Header returns the response headers.
func (this *LimitResponse) Header() http.Header {
	return this.Writer.Header()
}

// Write always succeeds and writes to this.Body, if not nil.
func (this *LimitResponse) Write(buf []byte) (bytes int, err error) {
	if this.exceeded {
		return
	}

	if (Bytes(len(buf)) + this.written) > this.Limit {
		this.exceeded = true
		return
	}

	bytes, err = this.Writer.Write(buf)
	this.written += Bytes(bytes)
	return
}

// WriteHeader sets this.Code.
func (this *LimitResponse) WriteHeader(code int) {
	this.Writer.WriteHeader(code)
}
