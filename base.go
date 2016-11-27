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
