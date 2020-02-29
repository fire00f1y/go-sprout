// The gcs package implements a resource for Google Storage. A client will be create lazily on first usage.
package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"sync"
)

var (
	c *storage.Client

	storageClientUninitialized = errors.New("[gosprout] storage client not initialized")

	mu = &sync.Mutex{}

	GzipContentType              = "application/gzip"
	XZipContentType              = "application/x-zip-compressed"
	ZipContentType               = "application/zip"
	OctetStreamContentType       = "application/octet-stream"
	BinaryOctetStreamContentType = "binary/octet-stream"
	TextPlainContentType         = "text/plain"
	XmlContentType               = "text/xml"
	XshContentType               = "text/x-sh"
	CsvContentType               = "text/csv"
)

func initClient(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	var err error
	c, err = storage.NewClient(ctx)

	return err
}

// Resource is a resource in google storage. To detect a change, we will poll on an interval and compare
// the version from the gcs metadata by using the Metageneration number and Generation number.
//
// When providing an UpdateFunc for these updates, you should consider the content type. Some typical content
// types have been defined in this package.
type Resource struct {
	bucket             string
	prefix             string
	lastMetageneration int64
	lastGeneration     int64
	contentType        string
}
