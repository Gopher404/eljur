package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type Errs interface {
	Error() string
	Down() Errs
}

type errs struct {
	err string
}

func New(err string) Errs {
	return &errs{
		err: "'" + err + "'",
	}
}

func NewE(err error) Errs {
	return New(err.Error())
}

func (e *errs) Error() string {
	return e.err
}

func (e *errs) Down() Errs {
	_, f, l, _ := runtime.Caller(1)
	e.err = fmt.Sprintf("%s:%d > %s", f[strings.LastIndex(f, "/")+1:], l, e.err)
	return e
}
