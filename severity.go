package carte

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type severity struct {
	name   []byte
	writer io.Writer

	hook func([]byte)

	mux sync.Mutex
}

// Log writes the writer specified in the severity struct
func (s *severity) Err(e error, details ...Jsonable) {
	details = append(details, &Jable{
		Name:  "ERR",
		Value: e.Error(),
	})
	_, _ = log(s, details...)
}

// Log writes the writer specified in the severity struct
func (s *severity) Log(details ...Jsonable) {
	_, _ = log(s, details...)
}

// Rlog is the same as log, but returns the values from the call to Write()
//func Rlog(sev *severity, details ...Jsonable) (int, error) {
//	return log(sev, details...)
//}

// Log writes the writer specified in the severity struct
func (s *severity) Msg(msg string, details ...Jsonable) {
	details = append(details, &Jable{
		Name:  "MSG",
		Value: msg,
	})
	_, _ = log(s, details...)
}

//Log writes the writer specified in the severity struct
func (s *severity) Print(details ...string) {
	j := &Jable{
		Name:  "DTLS",
		Value: fmt.Sprintf("%v", details),
	}
	_, _ = log(s, j)
}

func NewSeverity(name []byte, w io.Writer, hook func([]byte)) *severity {
	return &severity{
		name:   name,
		writer: w,
		hook:   hook,
	}
}

// GetWriter gets the writer of a severity
func (s *severity) GetWriter() io.Writer {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.writer
}

// SetWriter set the io.Writer that is logged to
// Use nil to ignore writer output
func (s *severity) SetWriter(w io.Writer) {
	s.mux.Lock()
	s.writer = w
	s.mux.Unlock()
}

func (s *severity) SetHook(hook func([]byte)) {
	s.mux.Lock()
	s.hook = hook
	s.mux.Unlock()
}

func (s *severity) getHook() func([]byte) {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.hook
}

// Predefined severity levels
var (
	// StdOut
	Info = &severity{
		name:   []byte("INF"),
		writer: os.Stdout,
	}
	Debug = &severity{
		name:   []byte("DBG"),
		writer: os.Stdout,
	}

	// StdErr
	Warn = &severity{
		name:   []byte("WRN"),
		writer: os.Stderr,
	}
	Error = &severity{
		name:   []byte("ERR"),
		writer: os.Stderr,
	}
	Critical = &severity{
		name:   []byte("CRT"),
		writer: os.Stderr,
	}
	Panic = &severity{
		name:   []byte("PNC"),
		writer: os.Stderr,
	}

	allSeverities = []*severity{Info, Debug, Warn, Error, Critical}
)

// SetWriters sets every severity to a single io.Writer
func SetWriters(w io.Writer, severities ...severity) {
	// Set writers
	for _, s := range severities {
		s.SetWriter(w)
	}
}

func SetAllWriters(w io.Writer) {
	// Set all writers
	for _, s := range allSeverities {
		s.SetWriter(w)
	}
}

// SetWriters sets every severity to a single io.Writer
func SetHookFor(hook func([]byte), severities ...severity) {
	// Set hooks
	for _, s := range severities {
		s.SetHook(hook)
	}
}

func SetAllHooks(hook func([]byte)) {
	// Set all hooks
	for _, s := range allSeverities {
		s.SetHook(hook)
	}
}
