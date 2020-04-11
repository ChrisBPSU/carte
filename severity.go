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

func (s *severity) GetWriter() io.Writer {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.writer
}

func (s *severity) SetWriter(w io.Writer) error {
	if w == nil {
		return ErrWriterWasNil
	}

	s.mux.Lock()
	s.writer = w
	s.mux.Unlock()

	return nil
}

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
