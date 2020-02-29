package gosprout

import (
	"io"
	"log"
)

var (
	// This will be called on any internal error. In some cases,
	// an external error handler is provided to a function. If this is the case, the provided
	// error handler will supercede the default one. If the provided one is nil, the default
	// will be used.
	DefaultErrorHandler ErrorHandler = func(e error) {
		log.Printf("[gosprout] %v\n", e)
	}
)

// Override the default error behavior. Be default, this will use the standard log.Printf.
func SetErrorHandler(errorHandler ErrorHandler) {
	DefaultErrorHandler = errorHandler
}

// WriteUpdate is a simple UpdateFunction. This will take in a provided writer which will be written to with the data
// from the resource specified. Closing will be handleded by the caller.
func WriteUpdate(w io.Writer) UpdateFunction {
	return func(r io.Reader) {
		_, err := io.Copy(w, r)
		if err != nil && DefaultErrorHandler != nil {
			DefaultErrorHandler(err)
		}
	}
}
