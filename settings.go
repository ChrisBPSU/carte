package carte

import (
	"errors"
	"sync"
	"time"
)

// Errors
var (
	ErrLocationWasNil = errors.New("received a nil time.Location")
)

// Settings
var (
	// -1 denotes the entire length
	functionNameLength = -1
	timezone           = time.UTC
	dateFormat         = "2006-01-02T15:04:05 MST"

	mux sync.Mutex
)

func SetFunctionNameLength(fnl int) {
	mux.Lock()
	functionNameLength = fnl
	mux.Unlock()
}

func SetTimezone(tz *time.Location) error {
	if tz == nil {
		return ErrLocationWasNil
	}

	mux.Lock()
	timezone = tz
	mux.Unlock()

	return nil
}

func SetDateFormat(format string) {
	mux.Lock()
	dateFormat = format
	mux.Unlock()
}
