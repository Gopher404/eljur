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
	value string
}

func New(err string) Errs {
	return &errs{
		value: "'" + err + "'",
	}
}

func NewE(err error) Errs {
	return New(err.Error())
}

func (e *errs) Error() string {
	return e.value
}

func (e *errs) Down() Errs {
	_, f, l, _ := runtime.Caller(1)
	e.value = fmt.Sprintf("%s:%d > %s", f[strings.LastIndex(f, "/")+1:], l, e.value)
	return e
}
