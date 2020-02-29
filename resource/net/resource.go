// The net resource is not implemented yet. The intention for the future is to be able to handle any remote network resource.
// This could be "http://", "https://", "tcp://", "ftp://" for example.
package net

import (
	"context"
	"errors"
	"io"
	"net/url"
)

var (
	notImplementedError = errors.New("[gosprout] net resource not yet implemented")
)

type Resource struct {
	url url.URL
}

func (netResource Resource) Poll(ctx context.Context) (bool, error) {
	return false, notImplementedError
}

func (netResource Resource) Refresh(ctx context.Context, out io.Writer, errorHandler func(error)) {
}
