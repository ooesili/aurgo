package logging

import (
	"fmt"
	"io"
)

func newLogger(out io.Writer) logger {
	return logger{out: out}
}

type logger struct {
	out io.Writer
}

func (l logger) log(format string, args ...interface{}) {
	format = "---> " + format + "\n"
	fmt.Fprintf(l.out, format, args...)
}
