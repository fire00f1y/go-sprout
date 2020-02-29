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

// This implements the Resource interface for an object in memory, for testing.
type MemTest struct {
	version      int
	data         string
	pollError    error
	handlerError error
}

func (m *MemTest) Poll(ctx context.Context) (bool, error) {
	if m.pollError != nil {
		return false, m.pollError
	}
	return true, nil
}

func (m *MemTest) Refresh(ctx context.Context, updateFunc func(io.Reader), errorFunc func(error)) {
	if m.handlerError != nil {
		errorFunc(m.handlerError)
		return
	}
	updateFunc(strings.NewReader(m.data))
}

func TestWatch_RefreshError(t *testing.T) {
	tests := []struct {
		interval     time.Duration
		data         string
		pollError    error
		handlerError error
	}{
		{
			interval:     500 * time.Millisecond,
			data:         "test test",
			pollError:    nil,
			handlerError: io.EOF,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			value := ""
			errorChan := make(chan error)
			ch := Watch(ctx,
				test.interval,
				&MemTest{
					version:      0,
					data:         test.data,
					pollError:    test.pollError,
					handlerError: test.handlerError,
				},
				func(r io.Reader) {
					t.Errorf("expected no update call, but got one\n")
				}, func(e error) {
					errorChan <- e
				})

			if value == test.data {
				t.Errorf("value was not supposed to be updated yet, but was. value: %s\n", value)
			}
			select {
			case <-time.After(test.interval + 100*time.Millisecond):
				{
					t.Errorf("timeout hit before error received")
				}
			case e := <-ch:
				{
					t.Errorf("error returned from watch channel unexpectedly: %v\n", e)
				}
			case <-errorChan:
				{
					log.Printf("test success\n")
				}
			}
		})
	}
}

func TestWatch_Error(t *testing.T) {
	tests := []struct {
		interval     time.Duration
		data         string
		pollError    error
		handlerError error
	}{
		{
			interval:     500 * time.Millisecond,
			data:         "test test",
			pollError:    io.EOF,
			handlerError: nil,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			value := ""
			ch := Watch(ctx,
				test.interval,
				&MemTest{
					version:      0,
					data:         test.data,
					pollError:    test.pollError,
					handlerError: test.handlerError,
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
		interval     time.Duration
		data         string
		pollError    error
		handlerError error
	}{
		{
			interval:     500 * time.Millisecond,
			data:         "test test",
			pollError:    nil,
			handlerError: nil,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			value := ""
			Watch(ctx,
				test.interval,
				&MemTest{
					version:      0,
					data:         test.data,
					pollError:    test.pollError,
					handlerError: test.handlerError,
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
