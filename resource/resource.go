// The resource package defines the api for a resource. It also provides a resource based on
// and input string. This is done via the scheme (the first part of the path provided).
// Schemes:
// - "gs://" will create a GCS resource
// - "file://" or "." or "/" or "\" (windows) will create a local file resource
// - "tcp://" or "http://" or "https://" or "ftp://" will create a network resource
//
// Custom resources can be defined by implementing the Resource interface defined in this package.
// For example, you could define one for a zookeeper resource which sets a ZK watch.
package resource

import (
	"context"
	"errors"
	"io"
)

var (
	UnknownTypeError    = errors.New("[gosprout] cannot derive resource type from path")
	NotImplementedError = errors.New("[gosprout] this is not implemented yet")
)

// The Resource needs to is a data source
type Resource interface {
	Poller
	Refresher
}

// Poller will find if the resource data source has been updated since the last time the resource was accessed.
type Poller interface {
	Poll(context.Context) (bool, error)
}

// Refreshers creates a reader from the underlying resource and calls the provided update func.
// The provided error handler will be used in case there is an issue.
type Refresher interface {
	Refresh(context.Context, func(io.Reader), func(error))
}

func CreateResource(path string) (Resource, error) {
	return nil, NotImplementedError
}
