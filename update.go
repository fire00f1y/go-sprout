package gosprout

import (
	"encoding/json"
	"io"
	"log"
	"sync"
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

type Serializer interface {
	Pointer() interface{}
	sync.Locker
}

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

// Provided a container object which is able to be locked, this function
// will lock, update the data in the underlying pointer, and unlock. The
// locking is necessary so this can be done without worry about concurrent access panics.
func UpdateFromJson(s Serializer) UpdateFunction {
	return func(r io.Reader) {
		p := s.Pointer()
		e := json.NewDecoder(r).Decode(p)
		if e != nil {
			log.Printf("[gosprout] failed to decode json object: %v\n", e)
		}
	}
}
