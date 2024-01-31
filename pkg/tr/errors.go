package tr

import (
	"fmt"
	"runtime"
)

func Trace(err error) error {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d > %w", tripPath(f), l, err) //f[strings.LastIndex(f, "/")+1:]
}

func tripPath(path string) string {
	const slash uint8 = 47
	s := false

	for i := len(path) - 1; i != 0; i-- {
		if path[i] == slash {
			if s {
				return path[i:]
			}
			s = true
		}
	}
	return path
}
