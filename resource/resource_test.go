package resource

import (
	"github.com/fire00f1y/go-sprout/resource/file"
	"github.com/fire00f1y/go-sprout/resource/gcs"
	"reflect"
	"strconv"
	"testing"
)

func TestCreateResource(t *testing.T) {
	tests := []struct {
		path          string
		typeStruct    interface{}
		expectedError error
	}{
		{
			path:          "gs://google-bucket/object",
			typeStruct:    gcs.Resource{},
			expectedError: nil,
		},
		{
			path:          "gs://google-bucket",
			typeStruct:    gcs.Resource{},
			expectedError: nil,
		},
		{
			path:          "file://google-bucket/object",
			typeStruct:    file.Resource{},
			expectedError: nil,
		},
		{
			path:          "./var/log",
			typeStruct:    file.Resource{},
			expectedError: nil,
		},
		{
			path:          "fake://var/log",
			typeStruct:    nil,
			expectedError: UnknownTypeError,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r, e := CreateResource(test.path)
			if e != test.expectedError {
				t.Errorf("expected error %v; got %v\n", test.expectedError, e)
			}

			if reflect.TypeOf(r) != reflect.TypeOf(test.typeStruct) {
				t.Errorf("type of created resource incorrect. expected %T; got %T\n", test.typeStruct, r)
			}
		})
	}

}

func TestGetScheme(t *testing.T) {
	tests := []struct {
		fullPath       string
		expectedScheme string
		expectedPath   string
		expectedError  error
	}{
		{
			fullPath:       "gs://google-bucket/test",
			expectedScheme: "gs",
			expectedPath:   "google-bucket/test",
			expectedError:  nil,
		},
		{
			fullPath:       "http://www.google.com",
			expectedScheme: "http",
			expectedPath:   "www.google.com",
			expectedError:  nil,
		},
		{
			fullPath:       "/var/log/java-app",
			expectedScheme: "file",
			expectedPath:   "/var/log/java-app",
			expectedError:  nil,
		},
		{
			fullPath:       "./www.google.com",
			expectedScheme: "file",
			expectedPath:   "./www.google.com",
			expectedError:  nil,
		},
		{
			fullPath:       ":/www.google.com",
			expectedScheme: "",
			expectedPath:   "",
			expectedError:  missingProtocolError,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s, p, e := getscheme(test.fullPath)
			if s != test.expectedScheme {
				t.Errorf("scheme failed: expected %s; got %s\n", test.expectedScheme, s)
			}
			if p != test.expectedPath {
				t.Errorf("path failed: expected %s; got %s\n", test.expectedPath, p)
			}
			if e != test.expectedError {
				t.Errorf("error failed: expected %v; got %v\n", test.expectedError, e)
			}
		})
	}
}
