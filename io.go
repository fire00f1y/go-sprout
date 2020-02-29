// gosprout is a library which is intended to watch resources and do some processing if the resource
// gets updated. A resource is defined in a subpackage. It is any source of data - a local file, a net resource,
//a GCS bucket, a database, etc...
package gosprout

import (
	"context"
	"github.com/fire00f1y/go-sprout/resource"
	"io"
	"time"
)

// UpdateFunction is an alias for a function which takes a reader. This is how
// a user of the library defines the functionality when the file has changed.
type UpdateFunction func(io.Reader)

// ErrorHandler is a function which takes an error. This will be called during
// processing if it encounters any issues.
type ErrorHandler func(error)

// Watch sets up a timer which will poll the resource on the provided interval and
// call the updateFunc if there is a new version available. The errorHandler will be
// called if there is any issue during processing. The provided ctx defines whether
// to continue or not - if it is Done() then updates will be permanently stopped.
//
// To use this, the user will need to provide its own logic of what to do with the data
// of a resource if it is a new version. There are some basic examples defined, but
// specific business logic will need to be provided in most cases.
func Watch(ctx context.Context,
	interval time.Duration,
	res resource.Resource,
	updateFunc UpdateFunction,
	errorHandler ErrorHandler) <-chan error {
	ch := make(chan error)

	go func(r resource.Resource) {
		timer := time.NewTimer(interval)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				{
					isNew, err := res.Poll(ctx)
					if err != nil {
						ch <- err
						continue
					}
					if isNew {
						r.Refresh(ctx, updateFunc, errorHandler)
					}
				}
			case <-ctx.Done():
				{
					close(ch)
					return
				}
			}
		}
	}(res)

	return ch
}
