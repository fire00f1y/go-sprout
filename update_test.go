package gosprout

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSetErrorHandler(t *testing.T) {
	f := func(c chan error) func(error) {
		return func(e error) {
			c <- e
		}
	}

	tests := []struct {
		ch       chan error
		expected error
	}{
		{
			ch:       make(chan error),
			expected: io.EOF,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			SetErrorHandler(f(test.ch))
			go DefaultErrorHandler(test.expected)
			select {
			case e := <-test.ch:
				{
					if e != test.expected {
						t.Errorf("expected %v; got %v\n", test.expected, e)
					}
				}
			case <-time.After(100 * time.Millisecond):
				{
					t.Error("no error provided after timeout")
				}
			}
		})
	}
}

func TestWriteUpdate(t *testing.T) {

	tests := []struct {
		text string
	}{
		{"test"},
		{"no real need for a second test"},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := new(bytes.Buffer)

			r := strings.NewReader(test.text)
			f := WriteUpdate(w)
			f(r)

			if w.String() != test.text {
				t.Errorf("expected %s; got %s\n", test.text, w.String())
			}
		})
	}
}
