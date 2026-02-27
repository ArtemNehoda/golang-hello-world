package logger

import (
	"log"
	"os"
)

// STDLogger wraps the standard library *log.Logger and satisfies Logger.
type STDLogger struct {
	l *log.Logger
}

// New returns a Logger backed by a standard library logger that writes to
// stderr with a date/time prefix.
func New() *STDLogger {
	return &STDLogger{
		l: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (s *STDLogger) Printf(format string, v ...any) { s.l.Printf(format, v...) }
func (s *STDLogger) Println(v ...any)               { s.l.Println(v...) }
func (s *STDLogger) Fatalf(format string, v ...any) { s.l.Fatalf(format, v...) }
func (s *STDLogger) Fatalln(v ...any)               { s.l.Fatalln(v...) }
