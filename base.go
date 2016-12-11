package goio

import "io"

type Flusher interface {
	Flush() error
}
type WriteFlushCloser interface {
	io.Writer
	Flusher
	io.Closer
}
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }


func NopCloser(w io.Writer) io.WriteCloser{
	return nopCloser{w}
}
