package file

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestResource_Poll(t *testing.T) {
	tests := []struct {
		path string
	}{
		{
			path: "test-file",
		}, {
			path: "./test-file2",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := Resource{
				path:       test.path,
				lastUpdate: time.Time{},
				lastSize:   0,
			}
			f, e := os.Create(test.path)
			if e != nil {
				t.Errorf("could not create test file: %v\n", e)
				return
			}
			defer os.Remove(test.path)
			defer f.Close()

			_, e = f.WriteString("test")
			if e != nil {
				t.Errorf("failed to write to test file: %v\n", e)
				return
			}
			updated, e := res.Poll(context.Background())
			if !updated {
				t.Errorf("file update did not successfully detect\n")
			}
			if e != nil {
				t.Errorf("error getting update from test file: %v\n", e)
			}
		})
	}
}

func TestResource_Refresh(t *testing.T) {
	tests := []struct {
		path string
		text string
	}{
		{
			path: "test-file",
			text: "test test test",
		}, {
			path: "./test-file2",
			text: "this is a test",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := Resource{
				path:       test.path,
				lastUpdate: time.Time{},
				lastSize:   0,
			}
			f, e := os.Create(test.path)
			if e != nil {
				t.Errorf("could not create test file: %v\n", e)
				return
			}
			defer os.Remove(test.path)

			_, e = f.WriteString(test.text)
			if e != nil {
				t.Errorf("failed to write to test file: %v\n", e)
				return
			}
			e = f.Close()
			if e != nil {
				t.Errorf("error closing file: %v\n", e)
				return
			}

			ch := make(chan string)
			go res.Refresh(context.Background(),
				func(r io.Reader) {
					b, _ := ioutil.ReadAll(r)
					ch <- string(b)
				}, func(e error) {
					t.Errorf("error durring file refresh: %v\n", e)
				})
			select {
			case s := <-ch:
				{
					if strings.TrimSpace(s) != strings.TrimSpace(test.text) {
						t.Errorf("expected text [%s]; got [%s]\n", test.text, s)
					}
				}
			case <-time.After(100 * time.Millisecond):
				{
					t.Error("timed out waiting for refresh response")
				}
			}
		})
	}
}
