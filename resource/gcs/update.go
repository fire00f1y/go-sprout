package gcs

import (
	"context"
	"io"
)

// Poll lazily initializes a storage client and then uses it to pull the attributes using the bucket and blob.
// The *ObjectAttrs which is returns contains metadata for the storage blob. We pull the ContentType, Metageneration,
// and Generation and store them locally to compare with the remote object. The ContentType should be referenced when
// deciding how to use the reader provided.
//
// See: https://pkg.go.dev/cloud.google.com/go/storage?tab=doc#ObjectAttrs
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

// Refresh currently does not distinguish between a file or a folder level object in GCS.
// It provides a reader for getting the data from the GCS object. It also manages the closing of the reader after
// completion.
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

	defer func() {
		if reader != nil {
			e := reader.Close()
			errorHandler(e)
		}
	}()
	if err != nil {
		go errorHandler(err)
		return
	}

	updateFunc(reader)
	return
}
