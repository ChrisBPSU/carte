package carte

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

// Errors
var (
	ErrNegativeFunctionNameLength = errors.New("received a negative function name length")
	ErrLocationWasNil             = errors.New("received a nil time.Location")
	ErrWriterWasNil               = errors.New("received a nil io.Writer")
)

type settings struct {
	functionNameLength int
	location           *time.Location
	dateFormat         string

	writerOut io.Writer
	writerErr io.Writer

	mux sync.Mutex
}

// Settings is a singleton
var (
	Settings = &settings{
		functionNameLength: 16,
		location:           time.UTC,
		dateFormat:         "2006-01-02T15:04:05 MST",
		writerOut:          os.Stdout,
		writerErr:          os.Stderr,
		mux:                sync.Mutex{},
	}
)

func (s *settings) SetFunctionNameLength(length int) error {
	if length < 0 {
		return ErrNegativeFunctionNameLength
	}

	s.mux.Lock()
	s.functionNameLength = length
	s.mux.Unlock()

	return nil
}

func (s *settings) SetLocation(loc *time.Location) error {
	if loc == nil {
		return ErrLocationWasNil
	}

	s.mux.Lock()
	s.location = loc
	s.mux.Unlock()

	return nil
}

func (s *settings) SetDateFormat(dateFormat string) {
	s.mux.Lock()
	s.dateFormat = dateFormat
	s.mux.Unlock()
}

func (s *settings) SetWriterOut(writerOut io.Writer) error {
	if writerOut == nil {
		return ErrWriterWasNil
	}

	s.mux.Lock()
	s.writerOut = writerOut
	s.mux.Unlock()

	return nil
}

func (s *settings) SetWriterErr(writerErr io.Writer) error {
	if writerErr == nil {
		return ErrWriterWasNil
	}

	s.mux.Lock()
	s.writerErr = writerErr
	s.mux.Unlock()

	return nil
}
