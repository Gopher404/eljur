package logger

import (
	"eljur/internal/config"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

func SetupLogger(cnf *config.LogConfig) (*slog.Logger, error) {
	var opts slog.HandlerOptions

	switch cnf.Level {
	case "debug":
		opts.Level = slog.LevelDebug
		break
	case "error":
		opts.Level = slog.LevelError
		break
	case "info":
		opts.Level = slog.LevelInfo
		break
	case "warn":
		opts.Level = slog.LevelWarn
		break
	default:
		return nil, fmt.Errorf("unknown level: %s", cnf.Level)
	}

	var handler slog.Handler

	var out io.Writer
	var outs []io.Writer
	for _, s := range strings.Split(cnf.Out, ",") {
		if s == "cmd" {
			outs = append(outs, os.Stdout)
		} else {
			f, err := os.Create(s)
			if err != nil {
				return nil, err
			}
			outs = append(outs, f)
			outFile = f
		}
	}
	out = io.MultiWriter(outs...)

	logType = cnf.Type

	switch cnf.Type {
	case "text":
		handler = slog.NewTextHandler(out, &opts)
		break
	case "json":
		handler = slog.NewJSONHandler(out, &opts)
		break
	default:
		return nil, fmt.Errorf("unknown type: %s", cnf.Type)
	}

	l := slog.New(handler)

	return l, nil
}

var logType string
var outFile *os.File
var mu sync.Mutex

func GetLogs() ([]byte, error) {
	if outFile == nil {
		return nil, errors.New("log file not init")
	}
	mu.Lock()
	_, _ = outFile.Seek(0, 0)
	logs, err := io.ReadAll(outFile)
	mu.Unlock()

	if err != nil {
		return nil, err
	}

	if logType != "json" {
		return logs, nil
	}

	logs = append([]byte("["), logs...)

	newLine := byte(10)
	l := len(logs)
	for i, b := range logs {
		if b == newLine && i+1 != l {
			logs[i] = byte(44) // add ","
		}
	}
	logs = append(logs, []byte("]")[0])
	return logs, nil
}
