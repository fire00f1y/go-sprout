package file

import (
	"os"
	"time"
)

// Resource is a local file. If the incoming path is in [".","file://","/"] it will be a
// file resource. The watching logic will not be a poller, but will instead be an event listener.
type Resource struct {
	path       string
	lastUpdate time.Time
	lastSize   int64
}

func NewResource(file string) (Resource, error) {
	info, err := os.Stat(file)
	if err != nil {
		return Resource{}, err
	}

	return Resource{
		path:       file,
		lastUpdate: info.ModTime(),
		lastSize:   info.Size(),
	}, err
}
