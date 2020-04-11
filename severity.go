package carte

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Errors
var (
	ErrWriterWasNil = errors.New("received a nil io.Writer")
)

type severity struct {
	name   []byte
	writer io.Writer

	mux sync.Mutex
}

// GetWriter gets the writer of a severity
func (s *severity) GetWriter() io.Writer {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.writer
}

// SetWriter set the io.Writer that is logged to
func (s *severity) SetWriter(w io.Writer) error {
	if w == nil {
		return ErrWriterWasNil
	}

	s.mux.Lock()
	s.writer = w
	s.mux.Unlock()

	return nil
}

// No nil checking, used internally for setting multiple writers
func (s *severity) setWriter(w io.Writer) {
	s.mux.Lock()
	s.writer = w
	s.mux.Unlock()
}

//Added ability to set all severities to a single writer, or to set multiple in a single call

// Predefined severity levels
var (
	// StdOut
	Info = &severity{
		name:   []byte("Info"),
		writer: os.Stdout,
	}
	Debug = &severity{
		name:   []byte("Debg"),
		writer: os.Stdout,
	}

	// StdErr
	Warn = &severity{
		name:   []byte("Warn"),
		writer: os.Stderr,
	}
	Error = &severity{
		name:   []byte("Err"),
		writer: os.Stderr,
	}
	Critical = &severity{
		name:   []byte("Crit"),
		writer: os.Stderr,
	}
)

// SetWriters sets every severity to a single io.Writer
func SetWriters(w io.Writer, severities ...severity) error {
	// Nil check writer
	if w == nil {
		return ErrWriterWasNil
	}

	// Set writers
	for _, s := range severities {
		s.setWriter(w)
	}

	return nil
}

func SetAllWriters(w io.Writer) error {
	// Nil check writer
	if w == nil {
		return ErrWriterWasNil
	}

	// Set all writers
	// On addition to new writers, add to here
	Info.setWriter(w)
	Debug.setWriter(w)
	Warn.setWriter(w)
	Error.setWriter(w)
	Critical.setWriter(w)

	return nil
}
