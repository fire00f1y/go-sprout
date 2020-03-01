package file

import (
	"context"
	"errors"
	"io"
	"os"
)

var (
	notImplementedError = errors.New("[gosprout] refresh from file not yet implemented")
)

func (r Resource) Poll(context.Context) (bool, error) {
	stat, err := os.Stat(r.path)
	if err != nil {
		return false, err
	}
	if stat.Size() != r.lastSize ||
		stat.ModTime() != r.lastUpdate {
		r.lastSize = stat.Size()
		r.lastUpdate = stat.ModTime()
		return true, nil
	}

	return false, nil
}

func (r Resource) Refresh(_ context.Context, updateFunc func(io.Reader), errorHandler func(error)) {
	f, err := os.Open(r.path)
	if err != nil {
		errorHandler(err)
		return
	}
	defer func() {
		e := f.Close()
		if e != nil {
			errorHandler(e)
		}
	}()

	updateFunc(f)
}
