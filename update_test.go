package gosprout

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"sync"
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

type container struct {
	data interface{}
	mu   *sync.Mutex
}

func (c container) Pointer() interface{} {
	return c.data
}

func (c container) Lock() {
	c.mu.Lock()
}

func (c container) Unlock() {
	c.mu.Unlock()
}

func TestUpdateFromJson(t *testing.T) {
	type jsonObject struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
	}

	tests := []struct {
		j    string
		name string
		n    int
	}{
		{j: `{"name":"test","number":1}`, name: "test", n: 1},
		{j: "this is not json, to make sure no panic happens", name: "", n: 0},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c := container{
				data: &jsonObject{},
				mu:   &sync.Mutex{},
			}

			r := strings.NewReader(test.j)
			f := UpdateFromJson(c)
			f(r)

			jj, ok := c.data.(*jsonObject)
			if !ok {
				t.Errorf("failed to type cast the container\n")
				return
			}

			if jj.Name != test.name || jj.Number != test.n {
				t.Errorf("decoded object does not match expected\n")
				return
			}
		})
	}
}
