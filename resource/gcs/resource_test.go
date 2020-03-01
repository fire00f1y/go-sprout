package gcs

import (
	"strconv"
	"testing"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		path   string
		bucket string
		blob   string
	}{
		{"google-bucket/one", "google-bucket", "one"},
		{"google-bucket/", "google-bucket", ""},
		{"google-bucket", "google-bucket", ""},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := NewResource(test.path)
			if err != nil {
				t.Errorf("error creating new resource: %v\n", err)
				return
			}

			if res.bucket != test.bucket {
				t.Errorf("bucket mismatch; expected: %s, got: %s\n", test.bucket, res.bucket)
			}

			if res.prefix != test.blob {
				t.Errorf("blob mismatch; expected: %s, got: %s\n", test.blob, res.prefix)
			}
		})
	}
}
