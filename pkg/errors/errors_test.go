package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func f2() Errs {
	return New("test error").Down()
}

func f1() Errs {
	return f2().Down()
}

func TestErrs(t *testing.T) {
	out := "errors_test.go:13 > errors_test.go:9 > 'test error'"
	err := f1()

	assert.Equal(t, out, err.Error())
}
