package gostd

import (
	"context"
	"fmt"
	"os"
)

type StringWriter struct {
	done    chan<- bool
	file    *os.File
	errfile *os.File
}

func (w *StringWriter) WriteLine(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	w.file.WriteString(s)
}

func (w *StringWriter) WriteError(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	w.errfile.WriteString(s)
}

func (w *StringWriter) Done() {
	w.done <- true
}

func NewStringWriter(ctx context.Context, name string, needErr bool) (*StringWriter, error) {
	file, err := os.OpenFile("files/"+name+".txt", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	var errfile *os.File = file
	if needErr {
		errfile, err = os.OpenFile("files/errors/"+name+"Err.txt", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			file.Close()
			return nil, err
		}
	}

	done := make(chan bool)
	sw := &StringWriter{done, file, errfile}

	go func() {
		defer file.Close()
		defer errfile.Close()

		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}()
	return sw, nil
}
