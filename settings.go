package carte

import (
	"errors"
	"sync"
	"time"
)

// Errors
var (
	ErrNegativeFunctionNameLength = errors.New("received a negative function name length")
	ErrLocationWasNil             = errors.New("received a nil time.Location")
)

// Settings
var (
	functionNameLength = 16
	timezone           = time.UTC
	dateFormat         = "2006-01-02T15:04:05 MST"

	mux sync.Mutex
)

func SetFunctionNameLength(fnl int) error {
	if fnl < 0 {
		return ErrNegativeFunctionNameLength
	}

	mux.Lock()
	functionNameLength = fnl
	mux.Unlock()

	return nil
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
