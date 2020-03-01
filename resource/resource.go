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
	"github.com/fire00f1y/go-sprout/resource/file"
	"github.com/fire00f1y/go-sprout/resource/gcs"
	"io"
	"os"
	"strings"
)

var (
	UnknownTypeError       = errors.New("[gosprout] cannot derive resource type from path")
	NotImplementedError    = errors.New("[gosprout] this is not implemented yet")
	missingProtocolError   = errors.New("[gosprout] missing protocol scheme")
	malformedProtocolError = errors.New("[gosprout] improper scheme format")
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
	s, p, e := getscheme(path)
	if e != nil {
		return nil, e
	}
	switch s {
	case "gs":
		{
			i := strings.Index(p, "/")
			return gcs.NewResource(path)
		}
	case "file":
		{
			return file.NewResource(p)
		}
	default:
		{
			return nil, UnknownTypeError
		}
	}
}

// This is shamelessly stolen from the standard library since it was not exported.
func getscheme(rawurl string) (scheme, path string, err error) {
	for i := 0; i < len(rawurl); i++ {
		c := rawurl[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-':
			if i == 0 {
				return "", stripPrecedingSlashes(rawurl), nil
			}
		case c == ':':
			if i == 0 {
				return "", "", missingProtocolError
			}
			return rawurl[:i], stripPrecedingSlashes(rawurl[i+1:]), nil
		case c == os.PathSeparator || c == '.':
			// This will only work if the it is not running on windows
			return "file", rawurl, nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", rawurl, nil
		}
	}
	return "", stripPrecedingSlashes(rawurl), nil
}

func stripPrecedingSlashes(s string) string {
	if strings.HasPrefix(s, "/") {
		return stripPrecedingSlashes(strings.TrimPrefix(s, "/"))
	}
	return s
}
