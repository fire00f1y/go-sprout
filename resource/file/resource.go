package file

import (
	"time"
)

// Resource is a local file. If the incoming path is in [".","file://","/"] it will be a
// file resource. The watching logic will not be a poller, but will instead be an event listener.
type Resource struct {
	path       string
	lastUpdate time.Time
	lastSize   int64
}

func NewResource(file string) Resource {
	return Resource{
		path:       file,
		lastUpdate: time.Now(),
		lastSize:   0,
	}
}
