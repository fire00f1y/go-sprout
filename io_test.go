package gosprout

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

type MemTest struct {
	version int
	data    string
	e       error
}

func (m MemTest) Poll(ctx context.Context) (bool, error) {
	if m.e != nil {
		return false, m.e
	}
	return true, nil
}

func (m MemTest) Refresh(ctx context.Context, updateFunc func(io.Reader), errorFunc func(error)) {
	updateFunc(strings.NewReader(m.data))
}

func TestWatch_Error(t *testing.T) {
	tests := []struct {
		interval time.Duration
		data     string
		testErr  error
	}{
		{
			interval: 500 * time.Millisecond,
			data:     "test test",
			testErr:  io.EOF,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			value := ""
			ch := Watch(ctx,
				test.interval,
				MemTest{
					version: 0,
					data:    test.data,
					e:       test.testErr,
				},
				func(r io.Reader) {
					t.Errorf("expected no update call, but got one\n")
				}, func(e error) {
					t.Errorf("expected no errors, but got %v\n", e)
				})

			if value == test.data {
				t.Errorf("value was not supposed to be updated yet, but was. value: %s\n", value)
			}
			select {
			case <-time.After(test.interval + 100*time.Millisecond):
				{
					t.Errorf("timeout hit before error received")
				}
			case <-ch:
				{
					log.Printf("test success\n")
				}
			}
		})
	}
}

func TestWatch_Update(t *testing.T) {

	tests := []struct {
		interval time.Duration
		data     string
	}{
		{
			interval: 500 * time.Millisecond,
			data:     "test test",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			value := ""
			Watch(ctx,
				test.interval,
				MemTest{
					version: 0,
					data:    test.data,
					e:       nil,
				},
				func(r io.Reader) {
					b, _ := ioutil.ReadAll(r)
					value = string(b)
				}, func(e error) {
					t.Errorf("internal test error: %v", e)
				})

			if value == test.data {
				t.Errorf("value was not supposed to be updated yet, but was. value: %s\n", value)
			}
			<-time.After(test.interval + 100*time.Millisecond)
			if value != test.data {
				t.Errorf("watch did not properly update. expected %s; got %s\n", test.data, value)
			}
		})
	}
}
