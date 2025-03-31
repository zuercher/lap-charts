package iox

import (
	"fmt"
	"io"
)

type Writer interface {
	io.Writer

	Print(args ...any) (int, error)
	Println(args ...any) (int, error)
	Printf(string, ...any) (int, error)
}

type writer struct {
	io.Writer
}

func NewWriter(w io.Writer) Writer {
	return &writer{w}
}
func (w *writer) Printf(format string, args ...any) (int, error) {
	s := fmt.Sprintf(format, args...)
	return w.Write([]byte(s))
}

func (w *writer) Println(args ...any) (int, error) {
	s := fmt.Sprintln(args...)
	return w.Write([]byte(s))
}

func (w *writer) Print(args ...any) (int, error) {
	s := fmt.Sprint(args...)
	return w.Write([]byte(s))
}
