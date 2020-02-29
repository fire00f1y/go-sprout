package gcs

import (
	"context"
	"io"
)

func (r Resource) Poll(ctx context.Context) (bool, error) {
	if c == nil {
		if e := initClient(ctx); e != nil {
			return false, e
		}
	}
	attrs, err := c.
		Bucket(r.bucket).
		Object(r.prefix).
		Attrs(ctx)

	if err != nil {
		return false, err
	}
	if attrs.Metageneration != r.lastMetageneration ||
		attrs.Generation != r.lastGeneration {
		r.contentType = attrs.ContentType
		r.lastMetageneration = attrs.Metageneration
		r.lastGeneration = attrs.Generation
		return true, nil
	}
	return false, nil
}

func (r Resource) Refresh(ctx context.Context, updateFunc func(io.Reader), errorHandler func(error)) {
	if c == nil {
		if e := initClient(ctx); e != nil {
			errorHandler(e)
			return
		}
	}
	reader, err := c.
		Bucket(r.bucket).
		Object(r.prefix).
		NewReader(ctx)

	if err != nil {
		go errorHandler(err)
		return
	}
	defer func() {
		e := reader.Close()
		errorHandler(e)
	}()

	updateFunc(reader)
	return
}
